package handler

import (
	"maple-bot/internal/bot/handler/command"
	"maple-bot/internal/bot/handler/event"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func RegisterHandlers(dispather *ext.Dispatcher) {
	dispather.AddHandler(handlers.NewCommand("start", command.Start))
	dispather.AddHandler(handlers.NewChatJoinRequest(nil, event.OnNewChatJoinRequest))
	// 方案一
	// 接受到聊天加入请求的时候，私聊发送认证链接（参数包含 chat_id user_id date（这个date不能明文，必须进行一个签名），同时批准进群，但是会去除该成员所有权限
	// 然后html 里面根据 chat_id 获取群验证的方式，启动相应的验证
	// 如果用户通过了验证，后台可以根据这个 chat_id 和 user_id 来进行成员权限的重新赋予 （会进行date的解签校验，看有没有超时，没有超时的话就解禁，通过验证但超时了话会 kick，然后让重新加群进行验证
}
