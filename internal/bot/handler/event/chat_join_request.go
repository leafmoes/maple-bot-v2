package event

import (
	"fmt"
	"maple-bot/internal/config"
	"maple-bot/internal/storage"
	"maple-bot/internal/util"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func OnNewChatJoinRequest(b *gotgbot.Bot, ctx *ext.Context) error {
	groupData, _ := storage.NewKVStorage().GetGroupData(ctx)
	process_mode := groupData.Config.ChatJoinRequest.ProcessMode
	remind_on_chat := groupData.Config.ChatJoinRequest.RemindOnChat
	userChatId := ctx.ChatJoinRequest.UserChatId
	userId := ctx.ChatJoinRequest.From.Id
	userFullName := fmt.Sprintf("%s %s", ctx.ChatJoinRequest.From.FirstName, ctx.ChatJoinRequest.From.LastName)
	groupChatId := ctx.ChatJoinRequest.Chat.Id
	groupChatTitle := ctx.ChatJoinRequest.Chat.Title
	date := ctx.ChatJoinRequest.Date
	signature, _ := util.EnPwdCode([]byte(fmt.Sprintf("%s%s%s", groupChatId, userId, date)))
	// process_mode 参数，同意进群并禁言 或者 私聊验证后同意进群
	// 从群配置里面获取模式，就这两种
	verifyUrl := fmt.Sprintf(
		"%s?action=%s&process_mode=%s&group_chat_id=%s&user_id=%s&signature=%s",
		config.App.TelegramBot.WebAppURL, "chat_join_request", process_mode, groupChatId, userId, signature)
	// 然后从配置中读取 开不开在群聊中提醒去私聊验证 的配置，来是否在群中提醒
	if remind_on_chat {
		_, err := ctx.ChatJoinRequest.Chat.SendMessage(
			b,
			fmt.Sprintf("%s正在请求加入群组，请前往私聊进行人机验证", userFullName), // 得创建一个at消息，还没写
			&gotgbot.SendMessageOpts{
				ParseMode: "HTML",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{Text: "前往私聊认证", Url: "https://t.me/" + b.Username},
						},
						{
							{Text: "通过", CallbackData: "approve_join_request"},
							{Text: "拒绝", CallbackData: "decline_join_request"},
							{Text: "拒绝并封禁", CallbackData: "decline_join_request_and_ban"},
						}},
				},
			})
		if err != nil {
			return fmt.Errorf("failed to send reminder message to group on new chat join request: %w", err)
		}
	}
	b.SendMessage(
		userChatId,
		fmt.Sprintf("你正在请求加入群聊 <pre>%s</pre> \n 请点击下方按钮验证你不是机器人。", groupChatTitle),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{Text: "验证", WebApp: &gotgbot.WebAppInfo{Url: verifyUrl}},
				}},
			},
		})
	// 如果是 AgreeJoinAndMute 就会批准，然后禁言，如果是 ArgeeJoinAfterVerify 就不进行处理
	if process_mode == "AgreeJoinAndMute" {
		_, err := ctx.ChatJoinRequest.Chat.ApproveJoinRequest(b, userId, &gotgbot.ApproveChatJoinRequestOpts{})
		if err != nil {
			return fmt.Errorf("failed to approve join request: %w", err)
		}
		ctx.ChatJoinRequest.Chat.SetPermissions(b, gotgbot.ChatPermissions{
			CanSendMessages: false,
		}, &gotgbot.SetChatPermissionsOpts{})
	}
	return nil
}
