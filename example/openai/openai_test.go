package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/sonic"

	"github.com/jun3372/uniai"
	"github.com/jun3372/uniai/client"
	"github.com/jun3372/uniai/request"
)

func Test_Completions(t *testing.T) {
	token := os.Getenv("UNIAI_API_OPENAI_TOKEN")
	if token == "" {
		t.Fatal("UNIAI_API_OPENAI_TOKEN is empty")
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
