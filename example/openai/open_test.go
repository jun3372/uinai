package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/sonic"

	uniai "github.com/jun3372/uinai"
	"github.com/jun3372/uinai/internal/client"
	"github.com/jun3372/uinai/request"
)

func Test_OpenAI(t *testing.T) {
	token := os.Getenv("UINAI_API_OPENAI_TOKEN")
	if token == "" {
		t.Fatal("UINAI_API_OPENAI_TOKEN is empty")
	}

	chat := uniai.New(
		client.WithType(client.OpenAI),
		client.WithHost("https://dashscope.aliyuncs.com/compatible-mode"),
		client.AddHeader("Authorization", token),
	)

	in := *request.NewRequest(
		request.WithModel("qwen2-72b-instruct"),
		request.WithTopP(0.9),
		request.WithStream(true),
		request.WithStop([]string{"im_end"}),
		request.WithMessages([]request.Messages{
			request.NewMessage(request.MessageRoleSystem, "你是一位现代诗人，能够轻松的写出李白和杜甫的风格的诗词"),
			request.NewMessage(request.MessageRoleUser, "请帮我写一首关于春天的诗"),
		}),
	)

	resp, err := chat.Completions(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}

	for item := range resp {
		bs, _ := sonic.ConfigDefault.MarshalToString(item)
		fmt.Println("resp.Item=", bs)
	}
}