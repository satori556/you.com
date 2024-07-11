package you

import (
	"context"
	"github.com/bincooo/emit.io"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	cookie    = "xxx"
	model     = CLAUDE_3_OPUS
	clearance = "xxx"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0"
)

func TestChat(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	chat := New(cookie, model, "http://127.0.0.1:7890")
	chat.LimitWithE(true)

	Exec("8081", "http://127.0.0.1:7890", os.Stdout, os.Stdout)
	defer Exit()

	messages := []Message{
		{
			Question: "你好",
			Answer:   "你好，有什么可以帮到您的？",
		},
	}

	query := `<|system|>
	<|end|>
	<|user|>
	现在几点了
	<|end|>
	
	你是一个智能机器人，你专注于选择工具的给用户使用的能力。
	你需要将任务拆解成多个子任务，并执行最后一个。 有时候，你可以依赖工具的运行结果，来更准确的回答 用户。
	工具使用了 JSON Schema 的格式声明，其中 toolId 是工具的 description 是工具的描述， parameters 是工具的参数，包括参数的类型和描述，required 是必填参数的列表。
	请你根据工具描述，决定回答问题或是使用工具。在完成任务过程中，USER 代表用户的输入， TOOL_RESPONSE代表工具运行结果。ASSISTANT 代表你的输出。
	你的每次输出都必须以 0,1 开头，代表是否需要调用工具:
	0: 不使用工具。
	
	1: 使用工具，返回工具调用的参数。 例如:
	USER: 你好呀 <|end|>
	ANSWER: 0: <|end|>
	USER: 今天杭州的天气如何 <|end|>
	ANSWER: 1: {"toolId":"testToolId","arguments":{"city": "杭州"}} <|end|> TOOL_RESPONSE: """
	晴天......
	"""
	USER: 今天杭州的天气适合去哪里玩? <|end|>
	ANSWER: 1: {"toolId":"testToolId2","arguments":{"query": "杭州 天气 去哪里玩"}} <| end|>
	TOOL_RESPONSE: """
	晴天. 西湖、灵隐寺、千岛湖......
	"""
	ANSWER: 0: <|end|>
	现在，我们开始吧!下面是你本次可以使用的工具:
	""" [
	{
	}, {
	"toolId": "tvy4P",
	"description": "获取用户当前时区的时间。", "parameters": {
	    "type": "object",
	    "properties": {
	    }
	},
	"required": []
	"""
	
	其中:工具 BrepI 已执行 下面是正式的对话内容:
	USER: 给QQ群组发送一条天气信息。 工具推荐:toolId = tvy4P ANSWER:`
	messages = append(messages, Message{
		Answer: query,
	})
	//query := "你是什么模型？"
	query = " "

	session, err := emit.NewDefaultSession("http://127.0.0.1:7890", nil, "127.0.0.1")
	if err != nil {
		logrus.Fatal(err)
	}

	if clearance == "" {
		response, e := emit.ClientBuilder(session).
			GET("http://127.0.0.1:8081/clearance").
			DoS(http.StatusOK)
		if e != nil {
			t.Fatal(e)
		}

		obj, e := emit.ToMap(response)
		if e != nil {
			t.Fatal(e)
		}

		data := obj["data"].(map[string]interface{})
		clearance = data["cookie"].(string)
		userAgent = data["userAgent"].(string)

		fileBytes, e := os.ReadFile("chat_test.go")
		if e != nil {
			t.Fatal(e)
		}
		fileContent := string(fileBytes)

		compile := regexp.MustCompile(`clearance = "([^"]?)"`)
		fileContent = compile.ReplaceAllString(fileContent, `clearance = "`+clearance+`"`)
		compile = regexp.MustCompile(`userAgent = "([^"]?)"`)
		fileContent = compile.ReplaceAllString(fileContent, `userAgent = "`+userAgent+`"`)
		_ = os.WriteFile("chat_test.go", []byte(fileContent), 0644)
	}

	chat.Client(session)
	chat.CloudFlare(clearance, userAgent, "")
	err = chat.Custom(context.Background(), "you/"+model, "xxx")
	if err != nil {
		t.Fatal(err)
	}

	ch, err := chat.Reply(context.Background(), messages, "[]", query)
	if err != nil {
		t.Fatal(err)
	}
	echo(ch, t)
}

func echo(ch chan string, t *testing.T) {
	content := ""
	for {
		message, ok := <-ch
		if !ok {
			break
		}

		if strings.HasPrefix(message, "error:") {
			t.Fatal(message)
		}

		if strings.HasPrefix(message, "limits:") {
			t.Log(message)
			continue
		}

		t.Log(message)
		content += message
	}

	t.Log(content)
}
