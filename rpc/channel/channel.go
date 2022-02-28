package channel

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
)

// different to io.Reader/io.Writer, this always ensure full msg is send/get
// if an error occoured in get, return nil
type RealiableChannel interface {
	Get() []byte
	Send([]byte) error
	Closed() bool
}

type Mux struct {
	baseChannel RealiableChannel
	SubChannels map[byte]RealiableChannel
	counter     byte
	getBuffer   map[byte][][]byte
}

func NewMux(baseChannel RealiableChannel) *Mux {
	return &Mux{
		baseChannel: baseChannel,
		SubChannels: make(map[byte]RealiableChannel),
		counter:     0,
		getBuffer:   make(map[byte][][]byte),
	}
}

type SubChannel struct {
	baseChannel     RealiableChannel
	idenficationKey byte
	mux             *Mux
	isClosed        bool
}

func (sc *SubChannel) Send(data []byte) error {
	return sc.baseChannel.Send(bytes.Join([][]byte{[]byte{sc.idenficationKey}, data}, []byte{}))
}

func (sc *SubChannel) Get() []byte {
	buf, ok := sc.mux.getBuffer[sc.idenficationKey]
	if ok {
		if len(buf) > 0 {
			sc.mux.getBuffer[sc.idenficationKey] = buf[1:]
			if len(buf[0]) == 0 {
				sc.isClosed = true
				return nil
			}
			return buf[0]
		}
	}
	for {
		data := sc.baseChannel.Get()
		if data == nil || len(data) == 1 {
			return nil
		}
		key := data[0]
		if key == sc.idenficationKey {
			if len(data) == 1 {
				sc.isClosed = true
				return nil
			}
			return data[1:]
		} else {
			sc.mux.getBuffer[key] = append(sc.mux.getBuffer[key], data[1:])
		}
	}
}

func (sc *SubChannel) Closed() bool {
	return sc.mux.baseChannel.Closed() || sc.isClosed
}

func (m *Mux) NewSubChannel() RealiableChannel {
	m.counter += 1
	sc := &SubChannel{
		baseChannel:     m.baseChannel,
		idenficationKey: m.counter,
		mux:             m,
	}
	m.SubChannels[m.counter] = sc
	return sc
}

func (m *Mux) GetSubChannel(key byte) RealiableChannel {
	sc, ok := m.SubChannels[key]
	if ok {
		return sc
	}
	if key > m.counter {
		m.counter = key - 1
	}
	return m.NewSubChannel()
}

type EncryptedChannel struct {
	Connect    RealiableChannel
	IsInitator bool
	encryptor  *EncryptionSession
	isClosed   bool
}

func (i *EncryptedChannel) Closed() bool {
	return i.isClosed
}

func (i *EncryptedChannel) initiateEncryptSession() error {
	// 发起者生成一对随机椭圆密钥
	initiatorPrivateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	// 生成盐
	salt := make([]byte, 16)
	rand.Read(salt)
	// 将公钥编码为二进制
	encodedInitiatorPublicKey, _ := x509.MarshalPKIXPublicKey(&initiatorPrivateKey.PublicKey)
	// 将公钥和盐合并，需要base64编码以避免出现多个 "-"
	keyForInit := base64.RawStdEncoding.EncodeToString(encodedInitiatorPublicKey) + "-" + base64.RawStdEncoding.EncodeToString(salt)

	// 将密钥信息发送给响应者
	err := i.Connect.Send([]byte(keyForInit))
	if err != nil {
		return err
	}
	// 从响应者接受响应者的公钥
	data := i.Connect.Get()
	if data == nil {
		return fmt.Errorf("Cannot Get PublicKey from Responder: %v", err)
	}
	responderPubKeyData, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return fmt.Errorf("error parsing public key: %v", err)
	}
	responderPublicKey := new(ecdsa.PublicKey)
	ecdsaKey, ok := responderPubKeyData.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("expected ECDSA public key, but got %v", responderPubKeyData)
	}
	*responderPublicKey = *ecdsaKey
	// 生成协商密钥
	i.encryptor = &EncryptionSession{PublicKey: responderPublicKey, PrivateKey: initiatorPrivateKey, Salt: salt}
	return i.encryptor.Init()
}

func (r *EncryptedChannel) waitForEncryptSession() error {
	initiatorPubkeyAndSaltData := r.Connect.Get()
	if initiatorPubkeyAndSaltData == nil {
		return fmt.Errorf("Cannot Get PublicKey from Initiator")
	}
	// 接收者生成一对随机椭圆密钥
	responderPrivateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	encodedResponderPublicKey, _ := x509.MarshalPKIXPublicKey(&responderPrivateKey.PublicKey)
	err := r.Connect.Send(encodedResponderPublicKey)
	if err != nil {
		return err
	}

	// 读取发起者发送的公钥和盐
	b64InitiatorPubkeyAndSalt := strings.Split(string(initiatorPubkeyAndSaltData), "-")
	if len(b64InitiatorPubkeyAndSalt) != 2 {
		return fmt.Errorf("incorrect initiator PublicKey and Salt Data format")
	}
	data, err := base64.StdEncoding.DecodeString(b64InitiatorPubkeyAndSalt[0])
	if err != nil {
		return fmt.Errorf("error 64decode public key: %v", err)
	}
	initiatorPublicKeyData, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return fmt.Errorf("error parsing public key: %v", err)
	}
	initiatorPublicKey := new(ecdsa.PublicKey)
	ecdsaKey, ok := initiatorPublicKeyData.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("expected ECDSA public key, but got %v", initiatorPublicKeyData)
	}
	*initiatorPublicKey = *ecdsaKey
	salt, err := base64.RawStdEncoding.DecodeString(b64InitiatorPubkeyAndSalt[1])
	if err != nil {
		return fmt.Errorf("error 64decode salt: %v", err)
	}

	// 生成协商密钥
	r.encryptor = &EncryptionSession{PublicKey: initiatorPublicKey, PrivateKey: responderPrivateKey, Salt: salt}
	return r.encryptor.Init()
}

func (e *EncryptedChannel) Send(data []byte) error {
	encyptedData := data[:]
	e.encryptor.Encrypt(encyptedData)
	return e.Connect.Send(encyptedData)
}

func (e *EncryptedChannel) Get() []byte {
	encyptedData := e.Connect.Get()
	//fmt.Println(string(encyptedData))
	if encyptedData == nil {
		e.isClosed = true
		return nil
	}
	e.encryptor.Decrypt(encyptedData)
	return encyptedData
}

func (e *EncryptedChannel) Init() error {
	if e.IsInitator {
		return e.initiateEncryptSession()
	} else {
		return e.waitForEncryptSession()
	}
}

type ChanChanel struct {
	ConnectW chan []byte
	ConnectR chan []byte
	IsClosed bool
}

func (c *ChanChanel) Closed() bool {
	return c.IsClosed
}

func (c *ChanChanel) Send(data []byte) (err error) {
	if data == nil {
		return fmt.Errorf("Send nil")
	}
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err = fmt.Errorf("%v", r)
		return
	}()
	c.ConnectW <- data
	return nil
}

func (c *ChanChanel) Get() []byte {
	r := <-c.ConnectR
	if r == nil {
		c.IsClosed = true
	}
	return r
}
