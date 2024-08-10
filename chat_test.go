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
	cookie    = "stytch_csrf_private_token=pl1PX2_e2E-oLtBCdlc2uKqU4n_WkwUwNBq_6u2huemE; you_subscription=freemium; ydc_stytch_session_jwt=eyJhbGciOiJSUzI1NiIsImtpZCI6Imp3ay1saXZlLTIwYzU3NzYzLTcxYzYtNDdmNC1hMWVmLTQ2NTZlMTU5ZTJlZSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsicHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCJdLCJhdXRoMF9pZCI6bnVsbCwiZXhwIjoxNzIyMzk3NDMwLCJodHRwczovL3N0eXRjaC5jb20vc2Vzc2lvbiI6eyJpZCI6InNlc3Npb24tbGl2ZS1iM2ZmNDk0OC1jMGNmLTQyNTctOGE1Ni0yMTczY2RhNGNjYWEiLCJzdGFydGVkX2F0IjoiMjAyNC0wNy0zMVQwMzozODo1MFoiLCJsYXN0X2FjY2Vzc2VkX2F0IjoiMjAyNC0wNy0zMVQwMzozODo1MFoiLCJleHBpcmVzX2F0IjoiMjAyNC0xMC0yOVQwMzozODo1MFoiLCJhdHRyaWJ1dGVzIjp7InVzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMF8xNV83KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTI1LjAuMC4wIFNhZmFyaS81MzcuMzYgRWRnLzEyNS4wLjAuMCIsImlwX2FkZHJlc3MiOiI5MS4xMDMuMTIxLjE1MCJ9LCJhdXRoZW50aWNhdGlvbl9mYWN0b3JzIjpbeyJ0eXBlIjoib3RwIiwiZGVsaXZlcnlfbWV0aG9kIjoiZW1haWwiLCJsYXN0X2F1dGhlbnRpY2F0ZWRfYXQiOiIyMDI0LTA3LTMxVDAzOjM4OjUwWiIsImVtYWlsX2ZhY3RvciI6eyJlbWFpbF9pZCI6ImVtYWlsLWxpdmUtZjZmMjhjYjItZjNlNi00Y2EyLWEwOTYtMmQ0YzYzNjEyOTViIiwiZW1haWxfYWRkcmVzcyI6InpoZWFhcHduQHNoYXJrbGFzZXJzLmNvbSJ9fV19LCJpYXQiOjE3MjIzOTcxMzAsImlzcyI6InN0eXRjaC5jb20vcHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCIsIm5iZiI6MTcyMjM5NzEzMCwic3ViIjoidXNlci1saXZlLWVkOTkwZTVhLTViZTktNDUxNS1iMjUwLWM2NDMxNDVmNjU0ZCJ9.jqSrei64gd7GvSs3X6qO_sAHBOr1-kjXT-CwyHDWG8Xk0ba-WKtOq_vnjZ3cHijXXrCF63nFTT4RVbDsUlQO-nOtvOKiP_Yq6_PMmhh4oZWBR4CeImUuRMUKWiukI8rAX90B1Z8UWAGNvuGZU0HnDRJ7EobK1HnADH-5wb0OUZnJLene3YoV1MpBH3kujtU9pEYOIczLB_l9eiSWbg6EsW78mmIyFfNXOxOh_9e4fA5g7JBAKt7Z0SjLIFp-QohyPsfprCrh13JtgI6L_dKcbtpoz0nx_kkfHZLKJZ9PTQHEfBVvG9-5ZEnVUg8WPGij0CSKxFVdous2I59-IwiglQ; ydc_stytch_session=VNBMOIK_uzDqCSeKMfyznN3bcjjtS3AX2CxWD8vq24a3; ab.storage.userId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3Auser-live-ed990e5a-5be9-4515-b250-c643145f654d%7Ce%3Aundefined%7Cc%3A1722397130561%7Cl%3A1722397130562; _clck=lsk3yg%7C2%7Cfnx%7C0%7C1673; ab.storage.deviceId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3Ac307a474-b7f2-2626-dbfc-115e47478c56%7Ce%3Aundefined%7Cc%3A1722397117519%7Cl%3A1722397130562; stytch_session=VNBMOIK_uzDqCSeKMfyznN3bcjjtS3AX2CxWD8vq24a3; _clsk=1aq1z26%7C1722397118865%7C1%7C0%7Cx.clarity.ms%2Fcollect; FPAU=1.2.1925151374.1722397119; daily_query_date=Wed%20Jul%2031%202024; _ga=GA1.1.524725508.1722397119; youpro_subscription=false; _gtmeec=eyJjb3VudHJ5IjoiNzgwMzI4NThiNTIwMDJlZTVkZTg2Nzk5ZjU3NjliY2NiNjk5YmQ2NjhkN2RlNTY2MDczZmVhM2IzMDhjODJmNiJ9; ab.storage.sessionId.dcee0642-d796-4a7b-9e56-a0108e133b07=g%3A709af877-093d-1c44-043b-24fb656f10d0%7Ce%3A1722398930565%7Cc%3A1722397130562%7Cl%3A1722397130565; FPID=FPID2.2.Sg7DRx24pl2nqKnY3UsPj3CUuSFfZZ%2BHFSowsqR5q9k%3D.1722397119; youchat_smart_learn=true; youchat_personalization=true; ai_model=gpt_4o; FPLC=DiiV17XIGrx2dW3bC7ypy92G5ATQvwQIfwonznz5epW75pRziPnqyabTVxqfNKSQ4HX836AzBRJ47SG3ucalp5FWdL0hkpkgeqFjZS20Pxxa2KGkCR%2FuTb79d9DqQA%3D%3D; uuid_guest_backup=55a82c4d-7248-4311-9638-578a08ac6ae5; FPGSID=1.1722397118.1722397118.G-WYGVQX1R23.ZGxKReYh-6HBZjUBwQIpMg; daily_query_count=0; _ga_2N7ZM9C56V=GS1.1.1722397118.1.1.1722397130.0.0.402548558; total_query_count=0; stytch_session_jwt=eyJhbGciOiJSUzI1NiIsImtpZCI6Imp3ay1saXZlLTIwYzU3NzYzLTcxYzYtNDdmNC1hMWVmLTQ2NTZlMTU5ZTJlZSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsicHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCJdLCJhdXRoMF9pZCI6bnVsbCwiZXhwIjoxNzIyMzk3NDMxLCJodHRwczovL3N0eXRjaC5jb20vc2Vzc2lvbiI6eyJpZCI6InNlc3Npb24tbGl2ZS1iM2ZmNDk0OC1jMGNmLTQyNTctOGE1Ni0yMTczY2RhNGNjYWEiLCJzdGFydGVkX2F0IjoiMjAyNC0wNy0zMVQwMzozODo1MFoiLCJsYXN0X2FjY2Vzc2VkX2F0IjoiMjAyNC0wNy0zMVQwMzozODo1MVoiLCJleHBpcmVzX2F0IjoiMjAyNC0xMC0yOVQwMzozODo1MFoiLCJhdHRyaWJ1dGVzIjp7InVzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMF8xNV83KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTI1LjAuMC4wIFNhZmFyaS81MzcuMzYgRWRnLzEyNS4wLjAuMCIsImlwX2FkZHJlc3MiOiI5MS4xMDMuMTIxLjE1MCJ9LCJhdXRoZW50aWNhdGlvbl9mYWN0b3JzIjpbeyJ0eXBlIjoib3RwIiwiZGVsaXZlcnlfbWV0aG9kIjoiZW1haWwiLCJsYXN0X2F1dGhlbnRpY2F0ZWRfYXQiOiIyMDI0LTA3LTMxVDAzOjM4OjUwWiIsImVtYWlsX2ZhY3RvciI6eyJlbWFpbF9pZCI6ImVtYWlsLWxpdmUtZjZmMjhjYjItZjNlNi00Y2EyLWEwOTYtMmQ0YzYzNjEyOTViIiwiZW1haWxfYWRkcmVzcyI6InpoZWFhcHduQHNoYXJrbGFzZXJzLmNvbSJ9fV19LCJpYXQiOjE3MjIzOTcxMzEsImlzcyI6InN0eXRjaC5jb20vcHJvamVjdC1saXZlLTlkZWE3ZGI1LTJlMTUtNGE3ZC1iYjFmLTJjYjg0ODBlYTliMCIsIm5iZiI6MTcyMjM5NzEzMSwic3ViIjoidXNlci1saXZlLWVkOTkwZTVhLTViZTktNDUxNS1iMjUwLWM2NDMxNDVmNjU0ZCIsInVzZXIiOnsiZW1haWwiOiJ6aGVhYXB3bkBzaGFya2xhc2Vycy5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmFtaWx5X25hbWUiOiIiLCJnaXZlbl9uYW1lIjoiIiwibmFtZSI6InpoZWFhcHduQHNoYXJrbGFzZXJzLmNvbSIsIm5pY2tuYW1lIjoiemhlYWFwd24iLCJzdHl0Y2hfdXNlcl9pZCI6InVzZXItbGl2ZS1lZDk5MGU1YS01YmU5LTQ1MTUtYjI1MC1jNjQzMTQ1ZjY1NGQiLCJzdWIiOiJ1c2VyLWxpdmUtZWQ5OTBlNWEtNWJlOS00NTE1LWIyNTAtYzY0MzE0NWY2NTRkIn19.ml31zxgvzTDRf2HXNBsqyOA2xjA8JO70ruftny3p3GBBWx59_B4uF4z10n14NbT8WRLao9-mGaOSRQAJLLxwCcG3vL3d3S_tdihxMGoWBXbQfe2SU0B-crkqDAFORhpNNmrQrcbdNzqZvh7wfYfLliXkT-npYG-ervTsMjtxfpX7ah7DZyOuLW6aDVnhj0ckPpcLjnV5IwRe7fFQjt1QykEkqxrP9GtUzbSVgUSWaVnWrEflhQzMqci75lXKExIEf5BJDrrAKmOMiWrjEa91QGeSGEXvw30ah5G_ZYKgF-yePML0aEdH0PrZhK3BFt2EGwdhgyY43QyX-zH8q9SBnQ; safesearch_guest=Moderate; uuid_guest=55a82c4d-7248-4311-9638-578a08ac6ae5; "
	model     = CLAUDE_3_OPUS
	clearance = "_cfuvid=tmAXz8I19LRJRTehBPtAK1UNshyLj6CRWIMQmq2AF34-1723259180396-0.0.1.1-604800000; __cf_bm=FY_DbtweHBaEzLib45KuyHlcmlN6isv5DNpnCENJwkk-1723259180-1.0.1.1-N3pxWUZXa447UDyG5dnjismHtTriQ400eBK_K._cRO1q2xTsFnr2pC6oDODflu4rwbH01KmvzAduDkbvxPFm_HUKdNsTJbN3lRF_nD4zrTk; "
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

	session, err := emit.NewJa3Session("http://127.0.0.1:7890", 120)
	if err != nil {
		logrus.Fatal(err)
	}

	if clearance == "" {
		response, e := emit.ClientBuilder(session).
			Ja3("yes").
			GET("http://127.0.0.1:8080/clearance").
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
	//err = chat.Custom(context.Background(), "you/"+model, "xxx")
	//if err != nil {
	//	t.Fatal(err)
	//}

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
