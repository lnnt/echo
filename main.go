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
	handler := dispatcher.NewEventDispatcher("", "")
	handler = handler.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		// 处理消息 event，这里简单打印消息的内容
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
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
