package model

type GroupData struct {
	GroupName string      `json:"group_name"`
	Config    GroupConfig `json:"config"`
}

type GroupConfig struct {
	WelcomeInfo     WelcomeInfoConfig     `json:"welcome_info"`
	ChatJoinRequest ChatJoinRequestConfig `json:"chat_join_request"`
	LogChannel      string                `json:"log_channel"`
	ExtraAdmins     []string              `json:"extra_admins"`
}

type WelcomeInfoConfig struct {
	WelcomeText string `json:"text"`       // 提供一些模板变量可以使用
	ParseMode   string `json:"parse_mode"` // 和TG的ParseMode一致
}

type ChatJoinRequestConfig struct {
	VerifyModes  []string `json:"verify_modes"`   // 成员加入需要进行的验证流程，比如 ["Github","Github-Star","ChiralCarbon"]
	ProcessMode  string   `json:"process_mode"`   // AgreeJoinAndMute   /  ArgeeJoinAfterVerify"
	RemindOnChat bool     `json:"remind_on_chat"` // 在群组中提醒去私聊认证（不建议开启，容易防洪
}
