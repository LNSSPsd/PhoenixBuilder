# BDump文件格式

> 感謝 [EillesWan](https://github.com/EillesWan) 提供的翻譯，本文檔內容主要來自其翻譯內容，存在部分改動。
>
> 部分註有 `譯註` 的內容均為此貢獻者所註。


> [Happy2018new](https://github.com/Happy2018new) & 修訂日誌<br/>
> 修訂 `容器` 相關<br/>
> 新增第 `39` 號操作的解析<br/>
> 同步 [66f2aa5](https://github.com/LNSSPsd/PhoenixBuilder/commit/66f2aa5b129e51a2154b64e5ff8bffc15290cf02) 中有關 `Bdump` 文件格式的更改


BDump v3 是個用於存儲*Minecraft*建築結構的文件格式。其內容由指示建造過程的命令組成。

按照一定的順序來寫下每一個方塊的ID的文件格式會因為包含空氣方塊而徒增文件大小，因此我們設計了一種新的文件格式，引入了「畫筆」，並讓一系列的指令控製其進行移動或放置方塊。

*\[註：畫筆絕非機器人的位置，而是一個引入的抽象的概念\]*










## 基本文件結構

BDump v3 文件的後綴名為`.bdx`，且文件頭為`BD@`, 代表本bdump文件已使用 brotli（PhoenixBuilder 使用的壓縮質量為`6`）進行壓縮。請註意，文件頭為`BDZ`的 BDump 文件同時存在，其使用 gzip 壓縮，然而包含這種文件頭的`.bdx`文件不被 PhoenixBuilder 支持，因為其棄用較早，目前難以再找到此類型的文件。我們將這種文件頭定義為「壓縮頭」(compression header)，並且在此壓縮頭後面的內容將以壓縮頭所表明的方式進行壓縮。

> 註: BDump v2 的文件後綴是 `.bdp`，且文件頭為 `BDMPS\0\x02\0`。

在壓縮頭之後的，壓縮後內容的起始字符為 `BDX\0`，且作者的遊戲名緊跟其後，並以 `\0` 表示其玩家名的表示完畢。*\[譯註：即若作者之遊戲名為Eilles，則此文件壓縮後的內容應以*`BDX\0Eilles\0`*開頭\]* 此後之文本即含參指令了，它們惜惜相依，緊緊相連。每個指令的ID占有1字節(byte)的空間，其正是`unsigned char`所占的空間。

所有的操作都基於一個用以標識「畫筆」所在位置的 `Vec3` 值。

*\[譯註：原諒我才疏學淺，冒昧在這裏註明一下：* `Vec3` *值指的是一個用以表示三維向量或坐標的值\]*

來來來，我們看看指令列表先。

數據類型定義如下：

* {整型}(int)：即全體整數集，可包含正整數、0、負整數
* {無符號整型}（亦稱非負整型）(unsigned int)：即全體非負整數集，可包含正整數和0
* `char`(單字)（亦稱字符）：一個1字節長的{整型}值
* `unsigned char`(無符號單字)（亦稱無符號字符或非負字符）：一個1字節長的{無符號整型}值
* `short`(短整)：一個2字節長的{整型}值
* `unsigned short`(無符號短整（亦稱非負短整）)：一個2字節長的{無符號整型}值
* `int32_t`：4字節長的{整型}數據
* `uint32_t`：4字節長的{無符號整型}數據
* `char *`：以`\0`(UTF-8編碼)結尾的字符串
* `int`：即`int32_t`
* `unsigned int`：即`uint32_t`
* `bool`(布爾)：1字節長的布爾(亦稱邏輯)數據，僅可為真(`true`, `1`)或假(`false`, `0`)

> 請註意：BDump文件的數字信息將會以<font style="color:red;">**大端字節序**</font>(big endian)又稱<font style="color:red;">**大端序**</font>記錄.
>
> 大小端字節序有何不同呢？
>
> *\[譯註：你完全可以去查百度、必應上面搜索出來的解析，那玩意肯定讓你半蒙半懂，但這玩意本身相對而言也並非十分絕對得重要，你看下面這個全蒙的也挺好。\]*
>
> 例如，一個`int32`的`1`在小端字節序的表示下，內存中是這樣的`01 00 00 00`，而大端為`00 00 00 01`。

*\[譯註：下面這表格中，我把調色板(palette)翻譯為了方塊池，純是因為意譯，但是，我也知道這樣失去了很多原文的趣味，我也在思索一種更好的翻譯……\]*

| ID                | 內部名                                     | 描述                                                         | 參數                                                         |
| ----------------- | ------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 1                 | `CreateConstantString`                     | 將特定的 `字符串` 放入 `方塊池` 。`字符串` 在 `方塊池` 中的 `ID` 將按照調用此命令的順序進行排序。如：你第一次調用這個命令的時候，對應 `字符串` 的 `ID` 為 `0` ，第二次就是 `1` 了。你最多只能添加到 `65535`<br/>*\[譯註：通常情況下，`字符串` 是一個方塊的 `英文ID名` ，如 `glass` \]* | `char *constantString` |
| 2                 | **已棄用且已移除**                          | - | - |
| 3                 | **已棄用且已移除**                          | - | - |
| 4                 | **已棄用且已移除**                          | - | - |
| 5                 | **已棄用且已移除**                          | - | - |
| 6                 | `AddInt16ZValue0`                          | 將畫筆的 `Z` 坐標增加 `value` | `unsigned short value` |
| 7                 | `PlaceBlock`                               | 在畫筆所在位置放置一個方塊。同時指定欲放置的方塊的 `數據值(附加值)` 為 `blockData` ，且該方塊在方塊池中的 `ID` 為 `blockConstantStringID` | `unsigned short blockConstantStringID`<br/>`unsigned short blockData` |
| 8                 | `AddZValue0`                               | 將畫筆的 `Z` 坐標增加 `1` | - |
| 9                 | `NOP`                                      | 擺爛，即不進行操作(`No Operation`) | - |
| 10, `0x0A`        | **已棄用且已移除**                          | - | - |
| 11, `0x0B`        | **已棄用且已移除**                          | - | - |
| 12, `0x0C`        | `AddInt32ZValue0`                          | 將畫筆的 `Z` 坐標增加 `value` | `unsigned int value` |
| 13, `0x0D`        | `PlaceBlockWithBlockStates`                | 在畫筆所在位置放置一個方塊。同時指定欲放置的方塊的 `方塊狀態` 為 `blockStatesString` ，且該方塊在方塊池中的 `ID` 為 `blockConstantStringID`<br/> `方塊狀態` 的格式形如 `["color":"orange"]` | `unsigned short blockConstantStringID`<br/>`char *blockStatesString` |
| 14, `0x0E`        | `AddXValue`                                | 將畫筆的 `X` 坐標增加 `1` | - |
| 15, `0x0F`        | `SubtractXValue`                           | 將畫筆的 `X` 坐標減少 `1` | - |
| 16, `0x10`        | `AddYValue`                                | 將畫筆的 `Y` 坐標增加 `1` | - |
| 17, `0x11`        | `SubtractYValue`                           | 將畫筆的 `Y` 坐標減少 `1` | - |
| 18, `0x12`        | `AddZValue`                                | 將畫筆的 `Z` 坐標增加 `1` | - |
| 19, `0x13`        | `SubtractZValue`                           | 將畫筆的 `Z` 坐標減少 `1` | - |
| 20, `0x14`        | `AddInt16XValue`                           | 將畫筆的 `X` 坐標增加 `value` 且 `value` 可正可負，亦或 `0` | `short value` |
| 21, `0x15`        | `AddInt32XValue`                           | 將畫筆的 `X` 坐標增加 `value`<br/>此指令與上一命令的不同點是此指令使用 `int32_t` 作為其參數 | `int value` |
| 22, `0x16`        | `AddInt16YValue`                           | 將畫筆的 `Y` 坐標增加 `value` （同上理） | `short value` |
| 23, `0x17`        | `AddInt32YValue`                           | 將畫筆的 `Y` 坐標增加 `value` （同上理） | `int value` |
| 24, `0x18`        | `AddInt16ZValue`                           | 將畫筆的 `Z` 坐標增加 `value` （同上理） | `short value` |
| 25, `0x19`        | `AddInt32ZValue`                           | 將畫筆的 `Z` 坐標增加 `value` （同上理） | `int value` |
| 26, `0x1A`        | `SetCommandBlockData`                      | **(推薦使用 `36` 號命令)** 在畫筆當前位置的方塊設置指令方塊的數據 *\[譯註：這裏可能是說，無論是啥方塊都可以加指令方塊的數據，但是嘞，只有指令方塊才能起效\]* | `unsigned int mode {脈沖=0, 重復=1, 連鎖=2}`<br/>`char *command`<br/>`char *customName`<br/>`char *lastOutput (此項無效，可被設為 '\0')`<br/>`int tickdelay`<br/>`bool executeOnFirstTick`<br/>`bool trackOutput`<br/>`bool conditional`<br/>`bool needsRedstone` |
| 27, `0x1B`        | `PlaceBlockWithCommandBlockData`           | **(推薦使用 `36` 號命令)** 在畫筆當前位置放置方塊池中 `ID` 為 `blockConstantStringID` 的方塊，且該方塊的 `方塊數據值(附加值)` 為 `blockData` 。放置完成後，為這個方塊設置 `命令方塊` 的數據(若可行的話) | `unsigned short blockConstantStringID`<br/>`unsigned short blockData`<br/>`unsigned int mode {脈沖=0, 重復=1, 連鎖=2}`<br/>`char *command`<br/>`char *customName`<br/>`char *lastOutput (此項無效，可被設為 '\0')`<br/>`int tickdelay`<br/>`bool executeOnFirstTick`<br/>`bool trackOutput`<br/>`bool conditional`<br/>`bool needRedstone` |
| 28, `0x1C`        | `AddInt8XValue`                            | 將畫筆的 `X` 坐標增加 `value`<br/>此指令與命令 `AddInt16XValue(20) `的不同點是此指令使用 `char` 作為其參數 | `char value //int8_t value` |
| 29, `0x1D`        | `AddInt8YValue`                            | 將畫筆的 `Y` 坐標增加 `value` （同上理） | `char value //int8_t value` |
| 30, `0x1E`        | `AddInt8ZValue`                            | 將畫筆的 `Z` 坐標增加 `value` （同上理） | `char value //int8_t value` |
| 31, `0x1F`        | `UseRuntimeIDPool`                         | 使用預設的 `運行時ID方塊池`<br/>`poolId`(預設ID) 是 PhoenixBuilder 內的值。網易MC( 1.17.0 @ 2.0.5 )下的 `poolId` 被我們定為 `117`。 每一個 `運行時ID` 都對應著一個方塊，而且包含其 `方塊數據值(附加值)`<br/>相關內容詳見 [PhoenixBuilder/resources](https://github.com/LNSSPsd/PhoenixBuilder/tree/main/resources)<br/>**已不再在新版本中被使用** | `unsigned char poolId` |
| 32, `0x20`        | `PlaceRuntimeBlock`                        | 使用特定的 `運行時ID` 在當前畫筆的位置放置方塊 | `unsigned short runtimeId`                                   |
| 33, `0x21`        | `placeBlockWithRuntimeId`                  | 使用特定的 `運行時ID` 在當前畫筆的位置放置方塊 | `unsigned int runtimeId`                                     |
| 34, `0x22`        | `PlaceRuntimeBlockWithCommandBlockData`    | 使用特定的 `運行時ID` 在當前畫筆的位置放置命令方塊，並設置其數據 | `unsigned short runtimeId`<br/>`unsigned int mode {脈沖=0, 重復=1, 連鎖=2}`<br/>`char *command`<br/>`char *customName`<br/>`char *lastOutput (此項無效，可被設為 '\0')`<br/>`int tickdelay`<br/>`bool executeOnFirstTick`<br/>`bool trackOutput`<br/>`bool conditional`<br/>`bool needRedstone` |
| 35, `0x23`        | `PlaceRuntimeBlockWithCommandBlockDataAndUint32RuntimeID` | 使用特定的 `運行時ID` 在當前畫筆的位置放置指令方塊，並設置其數據 | `unsigned int runtimeId`<br/>`unsigned int mode {脈沖 = 0, 循環 = 1, 連鎖 = 2}`<br/>`char *command`<br/>`char *customName`<br/>`char *lastOutput (此項無效，可被設為 '\0')`<br/>`int tickdelay`<br/>`bool executeOnFirstTick`<br/>`bool trackOutput`<br/>`bool conditional`<br/>`bool needRedstone` |
| 36, `0x24`        | `PlaceCommandBlockWithCommandBlockData`    | 根據給定的 `方塊數據值(附加值)` 在當前畫筆所在位置放置一個指令方塊，並設置其數據值 | `unsigned short data`<br/>`unsigned int mode {脈沖 = 0, 循環 = 1, 連鎖 = 2}`<br/>`char *command`<br/>`char *customName`<br/>`char *lastOutput (此項無效，可被設為 '\0')`<br/>`int tickdelay`<br/>`bool executeOnFirstTick`<br/>`bool trackOutput`<br/>`bool conditional`<br/>`bool needRedstone` |
| 37, `0x25`        | `PlaceRuntimeBlockWithChestData`           | 在畫筆所在位置放置一個 `runtimeId`(特定的 `運行時ID`) 所表示的方塊(如箱子、熔爐、唱片機等)，並向此方塊載入數據<br/>其中 `slotCount` 的數據類型為 `unsigned char`，因為我的世界用一個字節來存儲物品欄編號。此參數指的是要載入的次數，即要載入的 `ChestData` 結構體數量 | `unsigned short runtimeId` <br/> `unsigned char slotCount` <br/> `struct ChestData data` |
| 38, `0x26`        | `PlaceRuntimeBlockWithChestDataAndUint32RuntimeID` | 在畫筆所在位置放置一個 `runtimeId`(特定的 `運行時ID`) 所表示的方塊(如箱子、熔爐、唱片機等)，並向此方塊載入數據<br/>其中 `slotCount` 的數據類型為 `unsigned char`，因為我的世界用一個字節來存儲物品欄編號。此參數指的是要載入的次數，即要載入的 `ChestData` 結構體數量 | `unsigned int runtimeId`<br/>`unsigned char slotCount`<br/>`struct ChestData data` |
| 39, `0x27`        | `RecordBlockEntityData`                    | 記錄畫筆所在方塊的 `方塊實體` 數據，但亦可用於記錄其他信息<br/>`uint32_t length` 指代 `unsigned char buffer[length]` 的具體長度，而 `unsigned char buffer[length]` 自身則用於記錄信息<br/>應當說明的是，由於一些限製，`PhoenixBuilder` 在此處記錄的字段不是完整的 `NBT` | `uint32_t length`<br>`unsigned char buffer[length]` |
| 88, `'X'`, `0x58` | `Terminate`                                | 停止讀入。註意！雖然通常的結尾應該是 `XE` （2字節），但是用 `X` （1字節）是允許的 | - |
| 90, `0x5A`        | `isSigned` (此命令並非是一個真實的命令)      | 這是一個與其他命令功能稍有不同的命令，其參數應當出現在其前面，而這個指令呢也只能出現在文件的末尾。在不知道所以然的情況下，請不要使用它，因為無效的簽名會使得 `PhoenixBuilder` 無法去構建你的結構。詳見 `簽名` 部分。 | `unsigned char signatureSize` |

此表為 bdump v4 到 2022/1/29 為止的全部指令。

此外，對於 `struct ChestData` 數據結構，應當如下：

```
struct ChestData {
	char *itemName;
	unsigned char count;
	unsigned short data;
	unsigned char slotID;
}
```


（下述內容的其中一部分目前未被更新，除去部分已經棄用的命令外，其余應當正常運作）










## 文件樣例
下面是一些 `bdx` 文件的例子。
***

假設我們是一個熊孩子，來放置一個TNT在 `{3,5,6}`(**相對坐標**) 上，順帶地再放一個循環指令方塊，裏面寫著 `kill @e[type=tnt]` 還加了懸浮字 `Kill TNT!` ，且始終啟用，放在 `{3,6,6}` 上，再順手一點，我們放一塊惡臭的玻璃在 `{114514,15,1919810}` 上，一塊惡臭的鐵塊在 `{114514,15,1919800}` 上。好了，那麽未被壓縮的 BDX 文件應為如下：

`BDX\0DEMO\0\x01tnt\0\x1C\x03\x01repeating_command_block\0\x01glass\0\x01iron_block\0\x1E\x06\x1D\x05\x07\0\0\0\0\x10\x1B\0\x01\0\0\x01kill @e[type=tnt]\0Kill TNT!\0\0\0\0\0\0\x01\x01\0\0\x1D\x09\x19\0\x1D\x4B\x3C\x15\0\x01\xBF\x4F\x07\0\x02\0\0\x1E\xF6\x07\0\x03\0\0XE`

下面是偽代碼形式的指令表達法，便於我們觀察此結構具體的運作模式。

```assembly
author 'DEMO\0'
CreateConstantString 'tnt\0' ; 方塊ID: 0
AddInt8XValue 3 ; 畫筆位置: {3,0,0}
CreateConstantString 'repeating_command_block\0' ; 方塊ID: 1
CreateConstantString 'glass\0' ; 方塊ID: 2
CreateConstantString 'iron_block\0' ; 方塊ID: 3
AddInt8ZValue 6 ; 畫筆位置: {3,0,6}
AddInt8YValue 5 ; 畫筆位置: {3,5,6}
PlaceBlock (int16_t)0, (int16_t)0 ; TNT將會被放在 {3,5,6}
AddYValue ; *Y++, 畫筆位置: {3,6,6}
PlaceCommandBlockWithCommandBlockData (int16_t)1, (int16_t)0, 1, 'kill @e[type=tnt]\0', 'Kill TNT!\0', '\0', (int32_t)0, 1, 1, 0, 0 ; 指令方塊將會被放在 {3,6,6}
AddInt8YValue 9 ; 畫筆位置: {3,15,6}
AddInt32ZValue 1919804 ; 1919810: 00 1D 4B 3C = 01d4b3ch, 畫筆位置: {3,15,1919810}
AddInt32XValue 114511 ; 114511: 00 01 BF 4F = 01bf4fh, 畫筆位置: {114514,15,1919810}
PlaceBlock (int16_t)2,(int16_t)0 ; 玻璃將會被放在 {114514,15,1919810}
AddInt8ZValue -10 ; -10: F6 = 0f6h, 畫筆位置: {114514,15,1919800}
PlaceBlock (int16_t)3,(int16_t)0 ; 鐵塊 將會被放在 {114514,15,1919800}
Terminate
db 'E'
```
***
如果希望在畫筆所在位置放置一個 `正在燃燒的熔爐` ，且這個 `正在燃燒的熔爐` 的第一格和第三格分別是 `蘋果 * 3` 和 `鉆石 * 64` ，則那麽未被壓縮的 BDX 文件應為如下：

`BDX\x00DEMO\x00\x1f\x75\x26\x00\x00\x15\x2c\x02apple\x00\x03\x00\x00\x00diamond\x00\x40\x00\x00\x02XE`

下面是偽代碼形式的指令表達法，便於我們觀察此結構具體的運作模式。

```assembly
author 'DEMO\0' ; 設置作者為 'DEMO'
UseRuntimeIDPool (unsigned char)117 ; 117: 75
PlaceRuntimeBlockWithChestDataAndUint32RuntimeID (unsigned int)5420, (unsigned char)2 , 'apple\x00', (unsigned char)3, (unsigned short)0, (unsigned char)0, 'diamond\x00', (unsigned char)64, (unsigned short)0, (unsigned char)2
Terminate
db 'E'
```

以下是關於上述用到的 `PlaceRuntimeBlockWithChestDataAndUint32RuntimeID` 的相關解析。<br>
|參數|解釋|代碼片段|其他/備註|
|-|-|-|-|
|`PlaceRuntimeBlockWithChestDataAndUint32RuntimeID (unsigned int)5420`|在畫筆所在位置放置一個 `正在燃燒的熔爐`<br/>因為 `正在燃燒的熔爐` 在 `ID` 為 `117` 的 `運行時ID方塊池` 中的 `ID` 是 `5420` |`\x26\x00\x00\x15\x2c`|`5420` 在 `16` 進製下，其 `大端字節序` 表達為 `\x00\x00\x15\x2c`<br/>`unsigned int` 是 `正整數型` ，因此有 `4` 個字節|
|`(unsigned char)2`|向 `正在燃燒的熔爐` 載入 `2` 次數據(載入 `2` 個 `ChestData` 結構體)|`\x02`|`2` 在 `16` 進製下，其 `大端字節序` 表達為 `\x02`<br/>`unsigned char` 是 `無符號字節型` ，因此有 `1` 個字節|
|`apple\x00`|放入 `蘋果` |`apple\x00`|`char *` 是以 `\x00`(`UTF-8` 編碼)結尾的字符串|
|`(unsigned char)3`|`蘋果` 的數量為 `3`|`\x03`|`3` 在 `16` 進製下，其 `大端字節序` 表達為 `\x03`<br/>`unsigned char` 是 `無符號字節型` ，因此有 `1` 個字節|
|`(unsigned short)0`|`蘋果` 的 `物品數據值` 為 `0`|`\x00\x00`|`0` 在 `16` 進製下，其 `大端字節序` 表達為 `\x00\x00`<br/>`unsigned short` 是 `無符號短整型` ，因此有 `2` 個字節|
|`(unsigned char)0`|將 `蘋果` 放在第 `1` 個槽位|`\x00`|`0` 在 `16` 進製下，其 `大端字節序` 表達為 `\x00`<br/>`unsigned char` 是 `無符號字節型` ，因此有 `1` 個字節<br/>第一個槽位一般使用 `0` ，第二個槽位則為 `1` ，第三個槽位則為 `2` ，以此類推。|
|`diamond\x00`|放入 `鉆石`|`diamond\x00`|`char *` 是以 `\x00`(`UTF-8` 編碼)結尾的字符串|
|`(unsigned char)64`|`鉆石` 的數量為 `64`|`\x40`|`64` 在 `16` 進製下，其 `大端字節序` 表達為 `\x40`<br/>`unsigned char` 是 `無符號字節型` ，因此有 `1` 個字節|
|`(unsigned short)0`|`鉆石` 的 `物品數據值` 為 `0`|`\x00\x00`|`0` 在 `16` 進製下，其 `大端字節序` 表達為 `\x00\x00`<br/>`unsigned short` 是 `無符號短整型` ，因此有 `2` 個字節|
|`(unsigned char)2`|將 `鉆石` 放在第 `3` 個槽位|`\x02`|`2` 在 `16` 進製下，其 `大端字節序` 表達為 `\x02`<br/>`unsigned char` 是 `無符號字節型` ，因此有 `1` 個字節<br/>第一個槽位一般使用 `0` ，第二個槽位則為 `1` ，第三個槽位則為 `2` ，以此類推。|

您可以在 [PhoenixBuilder/resources](https://github.com/LNSSPsd/PhoenixBuilder/tree/main/resources) 查看 `運行時ID方塊池` 。<br>
本樣例采用的是 [PhoenixBuilder/resources/blockRuntimeIDs/netease/runtimeIds_117.json](https://github.com/LNSSPsd/PhoenixBuilder/blob/main/resources/blockRuntimeIDs/netease/runtimeIds_117.json) 所述之版本。










## 簽名
*PhoenixBuilder* 的 `0.3.5` 版本實現了一個 `bdump 文件簽名系統` ，用以辨認文件**真正的**發布者。

請註意， `bdx` 文件可不必被簽名，除非用戶打開了 `-S`（嚴格）開關。但這並不妨礙你去給他簽名，如果你為了簽名而簽名的話，則應確保其正常工作，因為 *PhoenixBuilder* 會拒絕處理簽名不正確的 `bdx` 文件。

我們使用基於 `RSA` 的哈希方法對 `BDX` 文件進行 `簽名` 。簽名時，相應的服務器會為每個用戶頒發一個單獨的認證集，然後 *PhoenixBuilder* 用相應的 `私鑰` 對文件進行 `簽名` ，並向對應的硬編碼服務器提供文件中根密鑰鏈接的 `公鑰` ，用於校驗 `BDX` 文件的真實發布者。

有關 `簽名` 的更多信息及詳細細節，另見 `fastbuilder/bdump/utils.go` : `SignBDXNew`/`VerifyBDXNew`