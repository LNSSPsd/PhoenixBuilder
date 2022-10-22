package I18n

var I18nDict_zh_CN map[uint16]string = map[uint16]string{
	ACME_FailedToGetCommand:             "未能读取ACME命令",
	ACME_FailedToSeek:                   "无效ACME文件，因为seek操作失败了。",
	ACME_StructureErrorNotice:           "文件结构错误",
	ACME_UnknownCommand:                 "未知ACME命令（文件错误）",
	Auth_BackendError:                   "后端错误",
	Auth_FailedToRequestEntry:           "未能请求租赁服入口，请检查租赁服等级设置是否关闭及租赁服密码是否正确。",
	Auth_HelperNotCreated:               "辅助用户尚未创建，请前往用户中心进行创建。",
	Auth_InvalidFBVersion:               "FastBuilder 版本无效，请更新。",
	Auth_InvalidHelperUsername:          "辅助用户的用户名无效，请前往用户中心进行设置。",
	Auth_InvalidToken:                   "无效Token，请重新登录。",
	Auth_InvalidUser:                    "无效用户，请重新登录。",
	Auth_ServerNotFound:                 "租赁服未找到，请检查租赁服是否对所有人开放。",
	Auth_UnauthorizedRentalServerNumber: "对应租赁服号尚未授权，请前往用户中心进行授权。",
	Auth_UserCombined:                   "该用户已经合并到另一个账户中，请使用新账户登录。",
	Auth_FailedToRequestEntry_TryAgain:  "未能请求租赁服入口，请稍后再试。",
	BDump_Author:                        "作者",
	BDump_EarlyEOFRightWhenOpening:      "未能读取文件，因为文件过早结束，可能已经损坏。",
	BDump_FailedToGetCmd1:               "未能获取到 cmd[pos:0] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd2:               "未能获取到 cmd[pos1] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd4:               "未能获取到 cmd[pos2] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd6:               "未能获取到 cmd[pos3] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd7_0:             "未能获取到 cmd[pos4] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd7_1:             "未能获取到 cmd[pos5] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd10:              "未能获取到 cmd[pos6] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd11:              "未能获取到 cmd[pos7] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetCmd12:              "未能获取到 cmd[pos8] 的任何参数，文件可能已经损坏",
	BDump_FailedToGetConstructCmd:       "未能获取到构建命令，文件可能已经损坏",
	BDump_FailedToReadAuthorInfo:        "未能读取作者信息，文件可能已损坏",
	BDump_FileNotSigned:                 "文件未签名",
	BDump_FileSigned:                    "文件已签名，持有者：%s",
	BDump_NotBDX_Invheader:              "不是bdx文件（无效文件头）",
	BDump_NotBDX_Invinnerheader:         "不是bdx文件（无效内部文件头）",
	BDump_SignedVerifying:               "文件已签名，正在验证...",
	BDump_VerificationFailedFor:         "因 %v 未能验证文件签名。",
	BDump_Warn_Reserved:                 "警告：BDump/Import：使用了保留字段\n",
	CommandNotFound:                     "未能找到此命令",
	ConnectionEstablished:               "成功连接到服务器。",
	Copyright_Notice_Bouldev:            "版权所有 (c) FastBuilder DevGroup, Bouldev 2022",
	Copyright_Notice_Contrib:            "代码参与: Ruphane, CAIMEO, CMA2401PT",
	Crashed_No_Connection:               "较长时间内未能建立连接",
	Crashed_OS_Windows:                  "按ENTER（回车）键来退出程序。",
	Crashed_StackDump_And_Error:         "Stack dump 于上方显示。错误为：",
	Crashed_Tip:                         "FastBuilder Phoenix 运行过程遇到问题",
	CurrentDefaultDelayMode:             "目前默认延迟模式",
	CurrentTasks:                        "任务列表：",
	DelayModeSet:                        "延迟模式已设定",
	DelayModeSet_DelayAuto:              "延迟值已自动设置为: %d",
	DelayModeSet_ThresholdAuto:          "延迟阈值已自动设置为: %d",
	DelaySet:                            "延迟已设定",
	DelaySetUnavailableUnderNoneMode:    "[delay set] 于 none 模式下不可用",
	DelayThreshold_OnlyDiscrete:         "延迟阈值只可在 discrete 模式下被设置。",
	DelayThreshold_Set:                  "延迟阈值已设置为 %d",
	ERRORStr:                            "错误",
	EnterPasswordForFBUC:                "请输入你的FastBuilder用户中心登录密码(不会显示): ",
	Enter_FBUC_Username:                 "输入你的FastBuilder用户中心用户名: ",
	Enter_Rental_Server_Code:            "请输入租赁服号: ",
	Enter_Rental_Server_Password:        "输入租赁服密码 (如果没有设置则直接按[Enter], 输入不会回显): ",
	ErrorIgnored:                        "已忽略此错误",
	Error_MapY_Exceed:                   "使用立体地图画时，MapY 应为范围在 20 至 255 的值（您输入的值为 %v)",
	FBUC_LoginFailed:                    "FastBuilder用户中心的用户名或密码无效",
	FBUC_Token_ErrOnCreate:              "创建Token文件时出错：",
	FBUC_Token_ErrOnGen:                 "生成临时Token时出错",
	FBUC_Token_ErrOnRemove:              "未能删除token文件: %v",
	FBUC_Token_ErrOnSave:                "保存Token时出错：",
	FileCorruptedError:                  "文件已被损坏",
	Get_Warning:                         "get命令已被set命令取代，未来将移除此命令，请迁移至set",
	IgnoredStr:                          "已忽略",
	InvalidFileError:                    "无效文件",
	InvalidPosition:                     "未获取到有效坐标。（可忽略此错误）",
	Lang_Config_ErrOnCreate:             "创建语言配置文件时出错：%v",
	Lang_Config_ErrOnSave:               "保存语言配置时出错：%v",
	LanguageName:                        "简体中文",
	LanguageUpdated:                     "语言偏好已更新",
	Logout_Done:                         "已从FastBuilder用户中心退出登录。",
	Menu_BackButton:                     "< 返回",
	Menu_Cancel:                         "取消",
	Menu_CurrentPath:                    "当前路径",
	Menu_ExcludeCommandsOption:          "排除命令方块内容",
	Menu_GetEndPos:                      "获取终点坐标",
	Menu_GetPos:                         "获取坐标",
	Menu_InvalidateCommandsOption:       "命令无效化",
	Menu_Quit:                           "退出程序",
	Menu_StrictModeOption:               "严格模式",
	NotAnACMEFile:                       "所提供的文件不是ACME建筑文件",
	Notice_CheckUpdate:                  "正在检查更新，请稍等…",
	Notice_iSH_Location_Service:         "您正在使用iSH模拟器，定位权限需要被用于保持后台运行，除此之外没有任何定位数据被记录或使用，您可以随时关闭它。",
	Notice_OK:                           "完成\n",
	Notice_UpdateAvailable:              "有新的PhoenixBuilder版本（%s）可用。\n",
	Notice_UpdateNotice:                 "请更新本软件。\n",
	Notice_ZLIB_CVE:                     "您的zlib版本（%s）包含已被证实的严重漏洞，我们建议您更新它，以免发生意外",
	Notify_NeedOp:                       "需要 OP 权限以正常工作。",
	Notify_TurnOnCmdFeedBack:            "需要 sendcommandfeedback 为 true，我们已经为你打开该选项，使用完后请按需关闭",
	Omega_WaitingForOP:                  "Omega系统正在等待 OP 权限",
	Omega_Enabled:                       "Omega系统已启用！",
	OpPrivilegeNotGrantedForOperation:   "未授予机器人 OP 权限，请授予 OP 权限后再进行此操作",
	Parsing_UnterminatedEscape:          "转义未终止",
	Parsing_UnterminatedQuotedString:    "字符串引号部分未终止",
	PositionGot:                         "已获得到起点坐标",
	PositionGot_End:                     "已获得终点坐标",
	PositionSet:                         "已设定坐标",
	PositionSet_End:                     "已设定终点坐标",
	QuitCorrectly:                       "正常退出",
	Sch_FailedToResolve:                 "未能解析文件",
	SelectLanguageOnConsole:             "请在控制台中选择新语言",
	ServerCodeTrans:                     "服务器号",
	SimpleParser_Int_ParsingFailed:      "解析器：未能处理整数形参数",
	SimpleParser_InvEnum:                "解析器：无效枚举值，可用值有：%s.",
	SimpleParser_Invalid_decider:        "解析器：无效决定子",
	SimpleParser_Too_few_args:           "解析器：参数过少",
	Special_Startup:                     "已启用语言：简体中文\n",
	/* Special SysError Translations, for innocent kids */
	SysError_EACCES:             "权限拒绝，请检查是否已经允许该程序访问对应文件。",
	SysError_EBUSY:              "文件被占用，请稍后再试。",
	SysError_EINVAL:             "无效文件输入。",
	SysError_EISDIR:             "输入文件为目录，无效输入。",
	SysError_ENOENT:             "对应文件不存在。",
	SysError_ETXTBSY:            "文件被占用，请稍后再试。",
	SysError_HasTranslation:     "对于 %s 的文件操作出错：%s",
	TaskCreated:                 "任务已创建",
	TaskDisplayModeSet:          "任务状态显示模式已经设置为: %s.",
	TaskFailedToParseCommand:    "未能解析命令: %v",
	TaskNotFoundMessage:         "未能根据所提供的任务ID找到有效任务。",
	TaskPausedNotice:            "[任务 %d] - 已暂停",
	TaskResumedNotice:           "[任务 %d] - 已恢复",
	TaskStateLine:               "ID %d - 命令行:\"%s\", 状态: %s, 延迟值: %d, 延迟模式: %s, 延迟阈值: %d",
	TaskStoppedNotice:           "[任务 %d] - 已停止",
	TaskTTeIuKoto:               "任务",
	TaskTotalCount:              "总数：%d",
	TaskTypeCalculating:         "正在计算",
	TaskTypeDied:                "已死亡",
	TaskTypePaused:              "已暂停",
	TaskTypeRunning:             "运行中",
	TaskTypeSpecialTaskBreaking: "特殊任务:正在终止",
	TaskTypeSwitchedTo:          "任务创建类型已经切换为：%s.",
	TaskTypeUnknown:             "未知",
	Task_D_NothingGenerated:     "[任务 %d] 无任何结构成功生成。",
	Task_DelaySet:               "[任务 %d] - 延迟已设置: %d",
	Task_ResumeBuildFrom:        "从第 %v 个方块处恢复构建",
	Task_SetDelay_Unavailable:   "[setdelay] 在 none 延迟模式下不可用",
	Task_Summary_1:              "[任务 %d] %v 个方块被更改。",
	Task_Summary_2:              "[任务 %d] 用时: %v 秒",
	Task_Summary_3:              "[任务 %d] 平均速度: %v 方块/秒",
	UnsupportedACMEVersion:      "不支持该版本ACME结构（仅支持acme 1.2文件版本）",
	Warning_UserHomeDir:         "警告 - 无法获取当前用户主目录，将设定homedir=\".\";\n",

	/* 菜单帮助部分的扩展 */

	Help_Tip_join_1:        "成功进入服务器,您可以在终端输入 ",
	Help_Tip_join_2:        " 获取帮助",
	Menu_Tip_MC_Command:    "执行我的世界原生指令,例如:.say 你好,世界",                   // 我的世界 指令在 FB 的执行
	Menu_Tip_FB_World_Chat: "FB 世界聊天",                                    // FB 世界聊天
	Menu_Tip_Exit:          "退出FB程序",                                     // FB 退出程序
	Menu_Tip_Help:          "FB 帮助菜单",                                    // FB 帮助菜单的帮助命令的额外扩展
	Menu_Tip_Lang:          "切换语言",                                       // FB 语言重新选择
	Menu_Tip_logout:        "从 FastBuilder 用户中心退出登录",                     // FB 退出登录
	Menu_Tip_Progress:      "设置是否显示进度条（显示建筑的进度百分比，方块总数，瞬时速度等信息。默认为true）", // FB 是否显示进度条
	Menu_Tip_Round:         "在指定坐标出为圆点画圆（很多参数,详情请使用 help round查询）",       // FB 画圆命令
	Menu_Tip_Get:           "获取当前机器人当前坐标并设置为FB导入建筑时的起点",                  // FB 获取当前机器人当前坐标并设置为FB导入建筑时的起点
	Menu_Tip_Set:           "设置导入建筑的起始点坐标",                               // FB 设置导入建筑的起始点坐标
	Menu_Tip_Task:          "有关 FB 正在进行的操作列表",                            // FB 任务命令
	Menu_Tip_Setend:        "设置导入建筑的终点坐标(不是必须设置)",                        // FB 设置导入建筑的终点坐标(不是必须设置)
	Menu_Tip_delay:         "设置发包方案(指令速度限制)",                             // FB 设置发包方案(指令速度限制)

	Menu_Tip_Cmd_MC_Command:    ".",        // 我的世界 指令在 FB 的执行
	Menu_Tip_Cmd_FB_World_Chat: ">",        // FB 世界聊天
	Menu_Tip_Cmd_Exit:          "exit",     // FB 退出程序
	Menu_Tip_Cmd_Help:          "help",     // FB 帮助菜单的帮助命令的额外扩展
	Menu_Tip_Cmd_Lang:          "lang",     // FB 语言重新选择
	Menu_Tip_Cmd_logout:        "exit",     // FB 退出登录
	Menu_Tip_Cmd_Progress:      "progress", // FB 是否显示进度条
	Menu_Tip_Cmd_Round:         "Round ",   // FB 画圆命令
	Menu_Tip_Cmd_Get:           "get",      // FB 获取当前机器人当前坐标并设置为FB导入建筑时的起点
	Menu_Tip_Cmd_Set:           "set",      // FB 设置导入建筑的起始点坐标
	Menu_Tip_Cmd_Task:          "task",     // FB 任务命令
	Menu_Tip_Cmd_Setend:        "setend",   // FB 设置导入建筑的终点坐标(不是必须设置)
	Menu_Tip_Cmd_delay:         "delay",    // FB 设置发包方案(指令速度限制)

	/*Help 部分的详细描述*/
	Help_No_Find: "未找到对应的 %s 帮助命令\n",

	Help_Help:     "help <页码:int> 打开帮助菜单的对应页码(当前版本仅支持第一页)\nhelp <命令:str> 详细描述对应命令的语法以及帮助,说明\n",
	Help_Exit:     "exit 退出应用程序(不需要任何参数) \n例:exit\n正常退出\n",
	Help_delay:    "delay 设置全局指令执行延迟和执行方案：\ndelay mode <delayMode:continuous|discrete|none>  设定默认的发包方案\ndelay threshold <threshold:int>  设定默认阈值（最大方块集合），仅在 discrete 方案下有效\ndelay set <Delay>  设定默认的发包延迟。在continuous模式下单位为微秒； 在discrete模式下单位为秒\n",
	Help_Lang:     "lang 切换语言\n例:请在控制台中选择新语言\n[1] English\n[2] English (UK)\n[3] 简体中文\n[4] 繁體中文（香港）\n[5] 繁體中文（台灣）\n(ID): 3\n语言偏好已更新\n(注意:此示范是在简体中文的情况下执行的命令)\n",
	Help_logout:   "logout 退出登录,在执行完成此命令后将会退出FB并且下次登录需要重新输入用户名和密码\n例:logout\n已从FastBuilder用户中心退出登录。\n正常退出\n",
	Help_Progress: "progress <value:bool> 设置是否显示进度条（显示建筑的进度百分比，方块总数，瞬时速度等信息。默认为true）\n例:progress True\n解析器：无效枚举值，可用值有：true, false.\nprogress true\n任务状态显示模式已经设置为: true.\nprogress false\n任务状态显示模式已经设置为: false.\n",
}
