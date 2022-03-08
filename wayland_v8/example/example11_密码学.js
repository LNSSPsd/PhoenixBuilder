// 暂时，我们只提供了一个简单的，基于 AES 的字符串加密和解密的实现
var rawData="这是要被加密的字符串"
var key="fastbuilder.js.v8.gamme"

var encryptedData=encryption.aesEncrypt(rawData,key)
engine.message("这是加密后的字符串"+encryptedData.cipherText)

var recoveredData=encryption.aesDecrypt(encryptedData.cipherText,key,encryptedData.iv)
engine.message("这是解密后的字符串"+recoveredData)

// 其实际 golang 实现为:
// 密码学实现：
// func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
//     padding := blockSize - len(ciphertext)%blockSize
//     padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//     return append(ciphertext, padtext...)
// }
//
// func PKCS7UnPadding(origData []byte) []byte {
//     length := len(origData)
//     unpadding := int(origData[length-1])
//     return origData[:(length - unpadding)]
// }
//AesEncrypt 加密函数
// func aesEncrypt(_plaintext, _key string) (string,string, error) {
//     plaintext:=[]byte(_plaintext)
//     key := []byte(_key)
//     key32:=make([]byte,32)
//     copy(key32,key)
//     c := make([]byte, aes.BlockSize+len(plaintext))
//     iv := c[:aes.BlockSize]
//
//     block, err := aes.NewCipher(key32)
//     if err != nil {
//         return "","", err
//     }
//     blockSize := block.BlockSize()
//     plaintext = PKCS7Padding(plaintext, blockSize)
//     blockMode := cipher.NewCBCEncrypter(block, iv)
//     crypted := make([]byte, len(plaintext))
//     blockMode.CryptBlocks(crypted, plaintext)
//     return hex.EncodeToString(crypted),hex.EncodeToString(iv), nil
// }
//
// // AesDecrypt 解密函数
// func aesDecrypt(_ciphertext, _key, _iv string) (string, error) {
//     ciphertext, _ :=hex.DecodeString(_ciphertext)
//     key := []byte(_key)
//     iv, _ :=hex.DecodeString(_iv)
//     key32:=make([]byte,32)
//     copy(key32,key)
//     block, err := aes.NewCipher(key32)
//     if err != nil {
//         return "", err
//     }
//     blockSize := block.BlockSize()
//     blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
//     origData := make([]byte, len(ciphertext))
//     blockMode.CryptBlocks(origData, ciphertext)
//     origData = PKCS7UnPadding(origData)
//     return string(origData), nil
// }

// 桥接器
// encryption encryption.aesEncrypt(text, key)
// encryption:=v8go.NewObjectTemplate(iso)
// global.Set("encryption", encryption)
// if err := encryption.Set("aesEncrypt",
//     v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
//     if text, ok := hasStrIn(info, 0, "encryption.aesEncrypt[text]"); !ok {
//         throwException("encryption.aesEncrypt", text)
//     } else {
//         if key, ok := hasStrIn(info, 1, "encryption.aesEncrypt[key]"); !ok {
//             throwException("encryption.aesEncrypt", key)
//         } else {
//             encryptOut,iv,err := aesEncrypt(text,key)
//             if err!=nil{
//                 throwException("encryption.aesEncrypt",err.Error())
//                 return nil
//             }else{
//                 result:=v8go.NewObjectTemplate(iso)
//                 jsEncryptOut, _ := v8go.NewValue(iso, encryptOut)
//                 jsIV, _ := v8go.NewValue(iso, iv)
//                 result.Set("cipherText",jsEncryptOut)
//                 result.Set("iv",jsIV)
//                 obj,_:=result.NewInstance(info.Context())
//                 return obj.Value
//             }
//         }
//     }
//     return nil
// }),
// ); err != nil {
//     panic(err)
// }
// if err := encryption.Set("aesDecrypt",
//     v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
//     if text, ok := hasStrIn(info, 0, "encryption.aesDecrypt[text]"); !ok {
//         throwException("encryption.aesDecrypt", text)
//     } else {
//         if key, ok := hasStrIn(info, 1, "encryption.aesDecrypt[key]"); !ok {
//             throwException("encryption.aesDecrypt", key)
//         } else {
//             if iv, ok := hasStrIn(info, 2, "encryption.aesDecrypt[iv]"); !ok {
//                 throwException("encryption.aesDecrypt", key)
//             } else{
//                 decryptOut,err := aesDecrypt(text,key,iv)
//                 if err!=nil{
//                     throwException("encryption.aesDecrypt",err.Error())
//                     return nil
//                 }else{
//                     value, _ := v8go.NewValue(iso, decryptOut)
//                     return value
//                 }
//             }
//         }
//     }
//     return nil
// }),
// ); err != nil {
//     panic(err)
// }