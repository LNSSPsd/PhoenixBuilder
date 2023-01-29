package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

type bookData struct {
	PageList []string
	Author   string
	Title    string
	Signed   bool
}

func parseBookData(book map[string]interface{}) (bookData, error) {
	var normal bool
	var got []interface{}
	var pageList []string = []string{}
	var author string = ""
	var title string = ""
	var signed bool = false
	// prepare
	_, ok := book["pages"]
	if ok {
		got, normal = book["pages"].([]interface{})
		if !normal {
			return bookData{}, fmt.Errorf("parseBookData: Crashed in book[\"pages\"]; book = %#v", book)
		}
		for key, value := range got {
			val, normal := value.(map[string]interface{})
			if !normal {
				return bookData{}, fmt.Errorf("parseBookData: Crashed in book[\"pages\"][\"%v\"]; book = %#v", key, book)
			}
			_, ok := val["text"]
			if ok {
				text, normal := val["text"].(string)
				if !normal {
					return bookData{}, fmt.Errorf("parseBookData: Crashed in book[\"pages\"][\"%v\"][\"text\"]; book = %#v", key, book)
				}
				pageList = append(pageList, text)
			}
		}
	}
	// pages
	_, ok = book["author"]
	if ok {
		author, normal = book["author"].(string)
		if !normal {
			return bookData{}, fmt.Errorf("parseBookData: Crashed in book[\"author\"]; book = %#v", book)
		}
	}
	// author
	_, ok = book["title"]
	if ok {
		title, normal = book["title"].(string)
		if !normal {
			return bookData{}, fmt.Errorf("parseBookData: Crashed in book[\"title\"]; book = %#v", book)
		}
	}
	// title
	_, signed = book["generation"]
	// generation
	return bookData{
		PageList: pageList,
		Author:   author,
		Title:    title,
		Signed:   signed,
	}, nil
}

func openBook(env *environment.PBEnvironment, bookRunTimeId int32) {
	env.Connection.(*minecraft.Conn).WritePacket(&packet.InventoryTransaction{
		TransactionData: &protocol.UseItemTransactionData{
			ActionType: 1,
			HeldItem: protocol.ItemInstance{
				Stack: protocol.ItemStack{
					ItemType: protocol.ItemType{
						NetworkID: bookRunTimeId,
					},
					Count:         1,
					CanBePlacedOn: []string{},
					CanBreak:      []string{},
					HasNetworkID:  false,
				},
			},
		},
	})
}

func writeText(env *environment.PBEnvironment, bookData bookData) {
	for key, value := range bookData.PageList {
		env.Connection.(*minecraft.Conn).WritePacket(&packet.BookEdit{
			Text:       value,
			PageNumber: byte(key),
		})
	}
}

func signBook(env *environment.PBEnvironment, bookData bookData) {
	if bookData.Signed {
		env.Connection.(*minecraft.Conn).WritePacket(&packet.BookEdit{
			ActionType: packet.BookActionSign,
			Title:      bookData.Title,
			Author:     bookData.Author,
		})
	}
}

// 总是打开 slot.horbar 0 处的书并写入文字，然后签名
func WriteTextToBook(Environment *environment.PBEnvironment, ItemData *types.ChestSlot) error {
	if CheckVersion() {
		bookRunTimeId := int32(ItemRunTimeID["writable_book"])
		openBook(Environment, bookRunTimeId)
		got, err := parseBookData(ItemData.ItemNBT)
		if err != nil {
			return fmt.Errorf("GetBookWithText: %v", err)
		}
		writeText(Environment, got)
		signBook(Environment, got)
		// 打开并写入数据，然后签名
		if got.Signed {
			Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
				EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
				NewItem: protocol.ItemInstance{
					Stack: protocol.ItemStack{
						ItemType: protocol.ItemType{
							NetworkID: int32(ItemRunTimeID["written_book"]),
						},
						Count:         1,
						CanBePlacedOn: ItemData.CanPlaceOn,
						CanBreak:      ItemData.CanDestroy,
						HasNetworkID:  false,
					},
				},
			})
		} else {
			Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
				EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
				NewItem: protocol.ItemInstance{
					Stack: protocol.ItemStack{
						ItemType: protocol.ItemType{
							NetworkID: int32(ItemRunTimeID["writable_book"]),
						},
						Count:         1,
						CanBePlacedOn: ItemData.CanPlaceOn,
						CanBreak:      ItemData.CanDestroy,
						HasNetworkID:  false,
					},
				},
			})
		}
		// 这只是为了告诉其他玩家发生了什么
		cmdsender := Environment.CommandSender.(*commands.CommandSender)
		cmdsender.SendWSCommandWithResponce("list")
		// 此举只是为了等待上述过程完成
	}
	return nil
}
