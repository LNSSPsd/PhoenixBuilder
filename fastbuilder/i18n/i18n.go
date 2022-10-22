package I18n

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LanguageEnglish_US         = "en_US"
	LanguageEnglish_UK         = "en_UK"
	LanguageSimplifiedChinese  = "zh_CN"
	LanguageTraditionalChinese = "zh_HK"
	LanguageTaiwanChinese      = "zh_TW"
)

var SelectedLanguage = LanguageEnglish_US

// !! IMPORTANT !! Please DO NOT change the order of items w/ prefix Auth_ !
// New items can be added anywhere AFTER Auth_ items !

const (
	Special_Startup = iota
	ACME_FailedToGetCommand
	ACME_FailedToSeek
	ACME_StructureErrorNotice
	ACME_UnknownCommand
	Auth_BackendError                   // 5
	Auth_FailedToRequestEntry           // 6
	Auth_HelperNotCreated               // 7
	Auth_InvalidFBVersion               // 8
	Auth_InvalidHelperUsername          // 9
	Auth_InvalidToken                   // 10
	Auth_InvalidUser                    // 11
	Auth_ServerNotFound                 // 12
	Auth_UnauthorizedRentalServerNumber // 13
	Auth_UserCombined                   // 14
	Auth_FailedToRequestEntry_TryAgain  // 15
	BDump_Author
	BDump_EarlyEOFRightWhenOpening
	BDump_FailedToGetCmd1
	BDump_FailedToGetCmd2
	BDump_FailedToGetCmd4
	BDump_FailedToGetCmd6
	BDump_FailedToGetCmd7_0
	BDump_FailedToGetCmd7_1
	BDump_FailedToGetCmd10
	BDump_FailedToGetCmd11
	BDump_FailedToGetCmd12
	BDump_FailedToGetConstructCmd
	BDump_FailedToReadAuthorInfo
	BDump_FileNotSigned
	BDump_FileSigned
	BDump_NotBDX_Invheader
	BDump_NotBDX_Invinnerheader
	BDump_SignedVerifying
	BDump_VerificationFailedFor
	BDump_Warn_Reserved
	CommandNotFound
	ConnectionEstablished
	Copyright_Notice_Bouldev
	Copyright_Notice_Contrib
	Crashed_No_Connection
	Crashed_OS_Windows
	Crashed_StackDump_And_Error
	Crashed_Tip
	CurrentDefaultDelayMode
	CurrentTasks
	DelayModeSet
	DelayModeSet_DelayAuto
	DelayModeSet_ThresholdAuto
	DelaySet
	DelaySetUnavailableUnderNoneMode
	DelayThreshold_OnlyDiscrete
	DelayThreshold_Set
	ERRORStr
	EnterPasswordForFBUC
	Enter_FBUC_Username
	Enter_Rental_Server_Code
	Enter_Rental_Server_Password
	ErrorIgnored
	Error_MapY_Exceed
	FBUC_LoginFailed
	FBUC_Token_ErrOnCreate
	FBUC_Token_ErrOnGen
	FBUC_Token_ErrOnRemove
	FBUC_Token_ErrOnSave
	FileCorruptedError
	Get_Warning
	IgnoredStr
	InvalidFileError
	InvalidPosition
	Lang_Config_ErrOnCreate
	Lang_Config_ErrOnSave
	LanguageName
	LanguageUpdated
	Logout_Done
	Menu_BackButton
	Menu_Cancel
	Menu_CurrentPath
	Menu_ExcludeCommandsOption
	Menu_GetEndPos
	Menu_GetPos
	Menu_InvalidateCommandsOption
	Menu_Quit
	Menu_StrictModeOption
	NotAnACMEFile
	Notice_CheckUpdate
	Notice_iSH_Location_Service
	Notice_OK
	Notice_UpdateAvailable
	Notice_UpdateNotice
	Notice_ZLIB_CVE
	Notify_NeedOp
	Notify_TurnOnCmdFeedBack
	Omega_Enabled
	Omega_WaitingForOP
	OpPrivilegeNotGrantedForOperation
	Parsing_UnterminatedEscape
	Parsing_UnterminatedQuotedString
	PositionGot
	PositionGot_End
	PositionSet
	PositionSet_End
	QuitCorrectly
	Sch_FailedToResolve
	SelectLanguageOnConsole
	ServerCodeTrans
	SimpleParser_Int_ParsingFailed
	SimpleParser_InvEnum
	SimpleParser_Invalid_decider
	SimpleParser_Too_few_args
	SysError_EACCES
	SysError_EBUSY
	SysError_EINVAL
	SysError_EISDIR
	SysError_ENOENT
	SysError_ETXTBSY
	SysError_HasTranslation
	SysError_NoTranslation // Do not add a translation for it!
	TaskCreated
	TaskDisplayModeSet
	TaskFailedToParseCommand
	TaskNotFoundMessage
	TaskPausedNotice
	TaskResumedNotice
	TaskStateLine
	TaskStoppedNotice
	TaskTTeIuKoto
	TaskTotalCount
	TaskTypeCalculating
	TaskTypeDied
	TaskTypePaused
	TaskTypeRunning
	TaskTypeSpecialTaskBreaking
	TaskTypeSwitchedTo
	TaskTypeUnknown
	Task_D_NothingGenerated
	Task_DelaySet
	Task_ResumeBuildFrom
	Task_SetDelay_Unavailable
	Task_Summary_1
	Task_Summary_2
	Task_Summary_3
	UnsupportedACMEVersion
	Warning_UserHomeDir

	/* 帮助菜单帮助部分的扩展 */
	Help_Tip_join_1
	Help_Tip_join_2

	Menu_Tip_MC_Command    // 我的世界 指令在 FB 的执行
	Menu_Tip_FB_World_Chat // FB 世界聊天
	Menu_Tip_Exit          // FB 退出程序
	Menu_Tip_Help          // FB 帮助菜单的帮助命令的额外扩展
	Menu_Tip_Lang          // FB 语言重新选择
	Menu_Tip_logout        // FB 退出登录
	Menu_Tip_Progress      // FB 是否显示进度条
	Menu_Tip_Round         // FB 画圆命令
	Menu_Tip_Get           // FB 获取当前机器人当前坐标并设置为FB导入建筑时的起点
	Menu_Tip_Set           // FB 设置导入建筑的起始点坐标
	Menu_Tip_Task          // FB 任务命令
	Menu_Tip_Setend        // FB 设置导入建筑的终点坐标(不是必须设置)
	Menu_Tip_delay         // FB 设置发包方案(指令速度限制)
	/* 帮助菜单的语法规则*/
	Menu_Tip_Cmd_MC_Command    // 我的世界 指令在 FB 的执行
	Menu_Tip_Cmd_FB_World_Chat // FB 世界聊天
	Menu_Tip_Cmd_Exit          // FB 退出程序
	Menu_Tip_Cmd_Help          // FB 帮助菜单的帮助命令的额外扩展
	Menu_Tip_Cmd_Lang          // FB 语言重新选择
	Menu_Tip_Cmd_logout        // FB 退出登录
	Menu_Tip_Cmd_Progress      // FB 是否显示进度条
	Menu_Tip_Cmd_Round         // FB 画圆命令
	Menu_Tip_Cmd_Get           // FB 获取当前机器人当前坐标并设置为FB导入建筑时的起点
	Menu_Tip_Cmd_Set           // FB 设置导入建筑的起始点坐标
	Menu_Tip_Cmd_Task          // FB 任务命令
	Menu_Tip_Cmd_Setend        // FB 设置导入建筑的终点坐标(不是必须设置)
	Menu_Tip_Cmd_delay         // FB 设置发包方案(指令速度限制)

	/* Help 命令的详细描述*/
	Help_Help
	Help_Exit
	Help_delay
	Help_Lang
	Help_logout
	Help_Progress
	Help_No_Find
)

var LangDict map[string]map[uint16]string = map[string]map[uint16]string{
	LanguageEnglish_US:         I18nDict_en_US,
	LanguageEnglish_UK:         I18nDict_en_UK,
	LanguageSimplifiedChinese:  I18nDict_zh_CN,
	LanguageTraditionalChinese: I18nDict_zh_HK,
	LanguageTaiwanChinese:      I18nDict_zh_TW,
}

var I18nDict map[uint16]string

func ShouldDisplaySpecial() bool {
	_, has := I18nDict[Special_Startup]
	return has
}

func HasTranslationFor(transtype uint16) bool {
	_, has := I18nDict[transtype]
	return has
}

func SelectLanguage() {
	config := loadConfigPath()
	curLangDict := make(map[uint16]string)
	{
		i := 1
		for lang := range LangDict {
			curLangDict[uint16(i)] = lang
			fmt.Printf("[%d] %s\n", i, LangDict[lang][LanguageName])
			i++
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("(ID): ") // No \n
		inp, _ := reader.ReadString('\n')
		inpl := strings.TrimRight(inp, "\r\n")
		parsedInt, err := strconv.Atoi(inpl)
		if err != nil {
			continue
		}
		if parsedInt <= 0 || parsedInt > len(curLangDict) {
			continue
		}
		SelectedLanguage = curLangDict[uint16(parsedInt)]
		break
	}
	if file, err := os.Create(config); err != nil {
		fmt.Println(T(Lang_Config_ErrOnCreate), err)
		fmt.Println(T(ErrorIgnored))
	} else {
		_, err = file.WriteString(SelectedLanguage)
		if err != nil {
			fmt.Println(T(Lang_Config_ErrOnSave), err)
			fmt.Println(T(ErrorIgnored))
		}
		file.Close()
	}
}

func UpdateLanguage() {
	langdict, aru := LangDict[SelectedLanguage]
	if !aru {
		panic("Updating to a language that not currently provided")
		return
	}
	I18nDict = langdict
	fmt.Printf("%s\n", T(LanguageUpdated))
}

func T(code uint16) string {
	r, has := I18nDict[code]
	if !has {
		r, has = I18nDict_en_US[code]
		if !has {
			return "???"
		}
	}
	return r
}

func loadConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("[PLUGIN] WARNING - Failed to obtain the user's home directory. made homedir=\".\";\n")
		homedir = "."
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder")
	os.MkdirAll(fbconfigdir, 0755)
	file := filepath.Join(fbconfigdir, "language")
	return file
}
