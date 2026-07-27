package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	// 注册消息处理器
	// 用于签名验证和消息解密，默认可以传递为空串。但如果你在开发者后台 > 事件与回调 > 加密策略中开启了加密，则必须传递 Encrypt Key 和 Verification Token
	handler := dispatcher.NewEventDispatcher("verificationToken", "eventEncryptKey")
	handler = handler.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		// 处理消息 event，这里简单打印消息的内容
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnCustomizedEvent("这里填入你要自定义订阅的 event 的 key，例如 out_approval", func(ctx context.Context, event *larkevent.EventReq) error {
		// 原生消息体
		fmt.Println(string(event.Body))
		fmt.Println(larkcore.Prettify(event.Header))
		fmt.Println(larkcore.Prettify(event.RequestURI))
		fmt.Println(event.RequestId())
		// 处理消息
		cipherEventJsonStr, err := handler.ParseReq(ctx, event)
		if err != nil {
			//  错误处理
			return err
		}
		plainEventJsonStr, err := handler.DecryptEvent(ctx, cipherEventJsonStr)
		if err != nil {
			//  错误处理
			return err
		}
		// 处理解密后的 消息体
		fmt.Println(plainEventJsonStr)
		return nil
	})
	// 注册 http 路由
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	// 启动 http 服务
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}
}
