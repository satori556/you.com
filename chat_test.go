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
	cookie    = "stytch_csrf_private_token=UlhJKn85_1dvKTp6eJ3mPW2GFkD0G3TtWo6tW7xq4ljK; you_subscription=freemium; ydc_stytch_session_jwt=eyJhbGciOiJSUzI1NiIsImtpZCI6Imp3ay1saXZlLTIwYzU3NzYzLTcxYzYtNDdmNC1hMWVmLTQ2NTZlMTU5ZTJlZSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsicHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCJdLCJhdXRoMF9pZCI6bnVsbCwiZXhwIjoxNzIyMzk4NTM3LCJodHRwczovL3N0eXRjaC5jb20vc2Vzc2lvbiI6eyJpZCI6InNlc3Npb24tbGl2ZS1kMGZmZTc3Ny0xNGU2LTQ0Y2ItYjZkOC1iYmJmZmMyY2MxZTUiLCJzdGFydGVkX2F0IjoiMjAyNC0wNy0zMVQwMzo1NzoxN1oiLCJsYXN0X2FjY2Vzc2VkX2F0IjoiMjAyNC0wNy0zMVQwMzo1NzoxN1oiLCJleHBpcmVzX2F0IjoiMjAyNC0xMC0yOVQwMzo1NzoxN1oiLCJhdHRyaWJ1dGVzIjp7InVzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMF8xNV83KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTI1LjAuMC4wIFNhZmFyaS81MzcuMzYgRWRnLzEyNS4wLjAuMCIsImlwX2FkZHJlc3MiOiI5MS4xMDMuMTIxLjE1NyJ9LCJhdXRoZW50aWNhdGlvbl9mYWN0b3JzIjpbeyJ0eXBlIjoib3RwIiwiZGVsaXZlcnlfbWV0aG9kIjoiZW1haWwiLCJsYXN0X2F1dGhlbnRpY2F0ZWRfYXQiOiIyMDI0LTA3LTMxVDAzOjU3OjE3WiIsImVtYWlsX2ZhY3RvciI6eyJlbWFpbF9pZCI6ImVtYWlsLWxpdmUtODgyYzIyMjQtNmZjOS00YmM5LWFmYWYtOWE4MTVkZjA1OTM1IiwiZW1haWxfYWRkcmVzcyI6InphZ3RoeHVxQHNoYXJrbGFzZXJzLmNvbSJ9fV19LCJpYXQiOjE3MjIzOTgyMzcsImlzcyI6InN0eXRjaC5jb20vcHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCIsIm5iZiI6MTcyMjM5ODIzNywic3ViIjoidXNlci1saXZlLTM5ZDMwYjQ4LTc5NTgtNDMxNS04ZjAzLTUwMzEwNGJiN2Y5YSJ9.ayoP0JUUv_yQMvdK1FYIMPZNlZ1ndaJ3RgfLF3kzA2nkz1jZsrDtu28oH8gsGigFp1t2UNyH1BBVnf-TduYXkE81mJRsGVMcGpRr49_5-xRGA28IOk87MCgNhghUxHXIvW4jK0fwZLa0clGPt6vJqHVWNUJ1_aUw7jADNNm-GWHnuATBXChHpJkvpzMv3Sl7OGbuhPcKVu6HbH2eYEpnea4O3q9jo3I4oi7PxoUvH22daT1o-fU-0xAUjjcKL8pPrCYW9-eKUGlPpeD4ce7wlX6OCsXgpw1DG84YbE_epmMrKrPuN0QVOWmlTUf-BA0pXmRQHOdSmvjJVdp2ylDBEg; ydc_stytch_session=wsCU7brSxdieYYJ4r6vkNSXmXFcvk0oLguh-GCPo7LOf; ab.storage.userId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3Auser-live-39d30b48-7958-4315-8f03-503104bb7f9a%7Ce%3Aundefined%7Cc%3A1722398237977%7Cl%3A1722398237978; _clck=qj3f8%7C2%7Cfnx%7C0%7C1673; ab.storage.deviceId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3Ad135baa5-fa65-63de-da5a-5de2438b70c7%7Ce%3Aundefined%7Cc%3A1722398220719%7Cl%3A1722398237978; stytch_session=wsCU7brSxdieYYJ4r6vkNSXmXFcvk0oLguh-GCPo7LOf; _clsk=1eqaxi5%7C1722398221803%7C1%7C0%7Cx.clarity.ms%2Fcollect; FPAU=1.2.2059970197.1722398222; daily_query_date=Wed%20Jul%2031%202024; _ga=GA1.1.1114695990.1722398221; youpro_subscription=false; _gtmeec=eyJjb3VudHJ5IjoiNzgwMzI4NThiNTIwMDJlZTVkZTg2Nzk5ZjU3NjliY2NiNjk5YmQ2NjhkN2RlNTY2MDczZmVhM2IzMDhjODJmNiJ9; ab.storage.sessionId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3A7df21712-c6ce-5aec-0f41-9f6848056703%7Ce%3A1722400037981%7Cc%3A1722398237977%7Cl%3A1722398237981; FPID=FPID2.2.UmxBpdKtA0RXRpW0kIc8PtAuUNBpsNeKOSgoQYA6JGA%3D.1722398221; youchat_smart_learn=true; youchat_personalization=true; ai_model=gpt_4o; FPLC=59BOO96xakN469LgZ3GOQJ7qLrffkXMxki9idg9SvsMP6LAxsRMfF4FkaYMW%2BmQxVD2g9UO18%2F9auIMRJ41fMlmNmfxXzuoMrhfiOKEQ9P%2BL9sphwty%2Fzh8z1FjOeA%3D%3D; uuid_guest_backup=c2b2a2c6-a4f9-488e-8cef-3240d65e6908; FPGSID=1.1722398221.1722398221.G-WYGVQX1R23.sgq9OA9C56FQS8hLDg6Kvg; daily_query_count=0; _ga_2N7ZM9C56V=GS1.1.1722398221.1.1.1722398237.0.0.1972730392; total_query_count=0; stytch_session_jwt=eyJhbGciOiJSUzI1NiIsImtpZCI6Imp3ay1saXZlLTIwYzU3NzYzLTcxYzYtNDdmNC1hMWVmLTQ2NTZlMTU5ZTJlZSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsicHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCJdLCJhdXRoMF9pZCI6bnVsbCwiZXhwIjoxNzIyMzk4NTM4LCJodHRwczovL3N0eXRjaC5jb20vc2Vzc2lvbiI6eyJpZCI6InNlc3Npb24tbGl2ZS1kMGZmZTc3Ny0xNGU2LTQ0Y2ItYjZkOC1iYmJmZmMyY2MxZTUiLCJzdGFydGVkX2F0IjoiMjAyNC0wNy0zMVQwMzo1NzoxN1oiLCJsYXN0X2FjY2Vzc2VkX2F0IjoiMjAyNC0wNy0zMVQwMzo1NzoxOFoiLCJleHBpcmVzX2F0IjoiMjAyNC0xMC0yOVQwMzo1NzoxN1oiLCJhdHRyaWJ1dGVzIjp7InVzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMF8xNV83KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTI1LjAuMC4wIFNhZmFyaS81MzcuMzYgRWRnLzEyNS4wLjAuMCIsImlwX2FkZHJlc3MiOiI5MS4xMDMuMTIxLjE1NyJ9LCJhdXRoZW50aWNhdGlvbl9mYWN0b3JzIjpbeyJ0eXBlIjoib3RwIiwiZGVsaXZlcnlfbWV0aG9kIjoiZW1haWwiLCJsYXN0X2F1dGhlbnRpY2F0ZWRfYXQiOiIyMDI0LTA3LTMxVDAzOjU3OjE3WiIsImVtYWlsX2ZhY3RvciI6eyJlbWFpbF9pZCI6ImVtYWlsLWxpdmUtODgyYzIyMjQtNmZjOS00YmM5LWFmYWYtOWE4MTVkZjA1OTM1IiwiZW1haWxfYWRkcmVzcyI6InphZ3RoeHVxQHNoYXJrbGFzZXJzLmNvbSJ9fV19LCJpYXQiOjE3MjIzOTgyMzgsImlzcyI6InN0eXRjaC5jb20vcHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCIsIm5iZiI6MTcyMjM5ODIzOCwic3ViIjoidXNlci1saXZlLTM5ZDMwYjQ4LTc5NTgtNDMxNS04ZjAzLTUwMzEwNGJiN2Y5YSIsInVzZXIiOnsiZW1haWwiOiJ6YWd0aHh1cUBzaGFya2xhc2Vycy5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmFtaWx5X25hbWUiOiIiLCJnaXZlbl9uYW1lIjoiIiwibmFtZSI6InphZ3RoeHVxQHNoYXJrbGFzZXJzLmNvbSIsIm5pY2tuYW1lIjoiemFndGh4dXEiLCJzdHl0Y2hfdXNlcl9pZCI6InVzZXItbGl2ZS0zOWQzMGI0OC03OTU4LTQzMTUtOGYwMy01MDMxMDRiYjdmOWEiLCJzdWIiOiJ1c2VyLWxpdmUtMzlkMzBiNDgtNzk1OC00MzE1LThmMDMtNTAzMTA0YmI3ZjlhIn19.TlZGqr2BsFD0IYlOwfkHOHdC_rpKWbJ0xLDmDWW-8bYW7OCjYtE-sqlv4fRsZRE99uxCgZwwBkBLmzE3NP-3FKUc0Mv-PkfXikjvQD1jh9Qg6nCGJG8sQSgMN1OV9Ei-LqZxuSYr7VU2Qd54MZEJ6UFIM4tkzedgkmm5ILhCjXiwZ-P1u4iIfeX4UfiIGseiRUYEsKNey6hWXu-9AG2Yj4GyTqBVUE0x9xzSTcOv1M5e4AM21U7aGuc_8F_3qYHLfta_t2g23wORDy4oGRo-GrSvomR9WKGXaSFXmQQI4hls51Sg6DMT-GibuU5QpIGO_TNsEsPGNwe5ciX4ZmZvYA; safesearch_guest=Moderate; uuid_guest=c2b2a2c6-a4f9-488e-8cef-3240d65e6908; "
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
