package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	client := lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"))
	handler := dispatcher.NewEventDispatcher(os.Getenv("TOKEN"), "")
	handler = handler.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Printf("[OnP2MessageReceiveV1 access], data: %s\n", larkcore.Prettify(event))

		var respContent map[string]string
		err := json.Unmarshal([]byte(*event.Event.Message.Content), &respContent)
		if err != nil || *event.Event.Message.MessageType != "text" {
			respContent = map[string]string{
				"text": "解析消息失败，请发送文本消息\nparse message failed, please send text message",
			}
		}

		content := larkim.NewTextMsgBuilder().
			TextLine("收到你发送的消息: " + respContent["text"]).
			TextLine("Received message: " + respContent["text"]).
			Build()

		if *event.Event.Message.ChatType == "p2p" {

			resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
				ReceiveIdType(larkim.ReceiveIdTypeChatId). // 消息接收者的 ID 类型，设置为会话ID。 ID type of the message receiver, set to chat ID.
				Body(larkim.NewCreateMessageReqBodyBuilder().
					MsgType(larkim.MsgTypeText).            // 设置消息类型为文本消息。 Set message type to text message.
					ReceiveId(*event.Event.Message.ChatId). // 消息接收者的 ID 为消息发送的会话ID。 ID of the message receiver is the chat ID of the message sending.
					Content(content).
					Build()).
				Build())

			if err != nil || !resp.Success() {
				fmt.Println(err)
				fmt.Println(resp.Code, resp.Msg, resp.RequestId())
				return nil
			}

		} else {
			resp, err := client.Im.Message.Reply(context.Background(), larkim.NewReplyMessageReqBuilder().
				MessageId(*event.Event.Message.MessageId).
				Body(larkim.NewReplyMessageReqBodyBuilder().
					MsgType(larkim.MsgTypeText). // 设置消息类型为文本消息。 Set message type to text message.
					Content(content).
					Build()).
				Build())
			if err != nil || !resp.Success() {
				fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
				return nil
			}
		}

		return nil
	})
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}
}
