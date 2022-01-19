package I18n

var I18nDict_en map[uint16]string = map[uint16]string{
	Copyright_Notice_Headline: "Copyright notice: \n",
	Copyright_Notice_Line_1: "FastBuilder Phoenix used codes\n",
	Copyright_Notice_Line_2: "from Sandertv's Gophertunnel that\n",
	Copyright_Notice_Line_3: "licensed under MIT license, at:\n",
	Crashed_Tip: "Oh no! FastBuilder Phoenix crashed!",
	Crashed_StackDump_And_Error: "Stack dump was shown above, error:",
	Crashed_OS_Windows: "Press ENTER to exit.",
	EnterPasswordForFBUC: "Enter your password for FastBuilder User Center: ",
	FBUC_LoginFailed: "Incorrect username or password",
	ServerCodeTrans: "Server",
	ConnectionEstablished: "Successfully created minecraft dialer.",
	InvalidPosition: "No position got. (ignorable)",
	PositionGot: "Position got",
	PositionGot_End: "End Position got",
	Enter_FBUC_Username: "Enter your FastBuilder User Center username: ",
	Enter_Rental_Server_Code: "Please enter your rental server number: ",
	Enter_Rental_Server_Password: "Enter Password (Press [Enter] if not set, input won't be echoed): ",
	NotAnACMEFile: "Invalid file, not an ACME structure.",
	UnsupportedACMEVersion: "Unsupported ACME structure version.Only acme file version 1.2 is supported.",
	ACME_FailedToSeek: "Invalid acme file since seek was failed.",
	ACME_FailedToGetCommand: "Failed to get acme command.",
	ACME_StructureErrorNotice: "Structure error",
	ACME_UnknownCommand: "Unknown ACME command",
	BDump_EarlyEOFRightWhenOpening: "Failed to read file, early EOF? File may be corrupted",
	BDump_NotBDX_Invheader: "Not a bdx file (Invalid file header)",
	InvalidFileError: "Invalid file",
	BDump_SignedVerifying: "File is signed, verifying...",
	FileCorruptedError: "File is corrupted",
	BDump_VerificationFailedFor: "Failed to verify the file's signature due to: %v",
	ERRORStr: "ERROR",
	IgnoredStr: "ignored",
	BDump_FileSigned: "File is signed, signer: %s",
	BDump_FileNotSigned: "File is not signed",
	BDump_NotBDX_Invinnerheader: "Not a bdx file (Invalid inner file header)",
	BDump_FailedToReadAuthorInfo: "Failed to read author info, file may be corrupted",
	BDump_Author: "Author",
	CommandNotFound: "Command not found.",
	Sch_FailedToResolve: "Failed to resolve file",
	SimpleParser_Too_few_args: "Parser: Too few arguments",
	SimpleParser_Invalid_decider: "Parser: Invalid decider",
	SimpleParser_Int_ParsingFailed: "Parser: failed to parse an int argument",
	SimpleParser_InvEnum: "Parser: Invalid enum value, allowed values are: %s.",
	QuitCorrectly: "Quit correctly",
	PositionSet: "Position set",
	PositionSet_End: "End position set",
	DelaySetUnavailableUnderNoneMode: "[delay set] is unavailable with delay mode: none",
	DelaySet: "Delay set",
	CurrentDefaultDelayMode: "Current default delay mode",
	DelayModeSet: "Delay mode set",
	DelayModeSet_DelayAuto: "Delay automatically set to: %d",
	DelayModeSet_ThresholdAuto: "Delay threshold automatically set to: %d",
	DelayThreshold_OnlyDiscrete: "Delay threshold is only available with delay mode: discrete",
	DelayThreshold_Set: "Delay threshold set to: %d",
	CurrentTasks: "Current tasks:",
	TaskStateLine: "ID %d - CommandLine:\"%s\", State: %s, Delay: %d, DelayMode: %s, DelayThreshold: %d",
	TaskTotalCount: "Total: %d",
	TaskNotFoundMessage: "Couldn't find a valid task by provided task id.",
	TaskPausedNotice: "[Task %d] - Paused",
	TaskResumedNotice: "[Task %d] - Resumed",
	TaskStoppedNotice: "[Task %d] - Stopped",
	Task_SetDelay_Unavailable: "[setdelay] is unavailable with delay mode: none",
	Task_DelaySet: "[Task %d] - Delay set: %d",
	TaskTTeIuKoto: "Task",
	TaskTypeSwitchedTo: "Task creation type set to: %s.",
	TaskDisplayModeSet: "Task status display mode set to: %s.",
	TaskCreated: "Task Created",
	Menu_GetPos: "getPos",
	Menu_GetEndPos: "getEndPos",
	Menu_Quit: "Quit Program",
	Menu_Cancel: "Cancel",
	Menu_ExcludeCommandsOption: "Exclude Commands",
	Menu_InvalidateCommandsOption: "Invalidate Commands",
	Menu_StrictModeOption: "Strict Mode",
	Menu_BackButton: "< Back",
	Menu_CurrentPath: "Current path",
	Parsing_UnterminatedQuotedString: "Unterminated quoted string",
	Parsing_UnterminatedEscape: "Unterminated escape",
	LanguageName: "English",
	TaskTypeUnknown: "Unknown",
	TaskTypeRunning: "Running",
	TaskTypePaused: "Paused",
	TaskTypeDied: "Died",
	TaskTypeCalculating: "Calculating",
	TaskTypeSpecialTaskBreaking: "SpecialTask:Breaking",
	TaskFailedToParseCommand: "Failed to parse command: %v",
	Task_D_NothingGenerated: "[Task %d] Nothing generated.",
	Task_Summary_1: "[Task %d] %v block(s) have been changed.",
	Task_Summary_2: "[Task %d] Time used: %v second(s)",
	Task_Summary_3: "[Task %d] Average speed: %v blocks/second",
	Logout_Done: "Logged out from FastBuilder User Center.",
	FailedToRemoveToken: "Failed to remove token file: %v",
	SelectLanguageOnConsole: "Please select your new preferred language on console.",
	LanguageUpdated: "Language preference has been updated",
	Auth_ServerNotFound: "Server not found, please check your server's public state",
	Auth_FailedToRequestEntry: "Failed to request entry for your server, please check whether the password is correct and please turn off the level limitation",
	Auth_InvalidHelperUsername: "Invalid username for helper user, please set it on FastBuilder User Center.",
	Auth_BackendError: "Backend Error",
	Auth_UnauthorizedRentalServerNumber: "Unauthorized rental server number, please add it on your FastBuilder User Center.",
	Auth_HelperNotCreated: "Helper user haven't been created, please go create it on FastBuilder User Center.",
	Auth_InvalidUser: "Invalid user for FastBuilder User Center",
	Auth_InvalidToken: "Invalid login token.",
	Auth_UserCombined: "Given user has been combined to another account, please login using another account's information.",
	Auth_InvalidFBVersion: "Invalid FastBuilder version, please update.",
	Notify_TurnOnCmdFeedBack:            "FastBuilder require sendcommandfeedback=true, please input:\"/gamerule sendcommandfeedfack true\"and restart FastBuilder.",
	Notify_NeedOp:                       "FastBuilder require OP privilege.",

}