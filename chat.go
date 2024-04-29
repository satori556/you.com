package you

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bincooo/you.com/common"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Role    string
	Content string
}

type Chat struct {
	cookie      string
	model       string
	proxies     string
	limitsWithE bool
}

const (
	GPT_4                    = "gpt_4"
	GPT_4_TURBO              = "gpt_4_turbo"
	CLAUDE_2                 = "claude_2"
	CLAUDE_3_OPUS            = "claude_3_opus"
	CLAUDE_3_SONNET          = "claude_3_sonnet"
	CLAUDE_3_HAIKU           = "claude_3_haiku"
	GEMINI_PRO               = "gemini_pro"
	GEMINI_1_5_PRO           = "gemini_1_5_pro"
	DATABRICKS_DBRX_INSTRUCT = "databricks_dbrx_instruct"
	COMMAND_R                = "command_r"
	COMMAND_R_PLUS           = "command_r_plus"
	LLAMA3                   = "llama3"
	ZEPHYR                   = "zephyr"
)

func New(cookie, model, proxies string) Chat {
	cookie = extCookies(cookie, model)
	return Chat{cookie, model, proxies, false}
}

func (c *Chat) Reply(ctx context.Context, previousMessages, query string) (chan string, error) {
	if c.limitsWithE {
		response, err := common.ClientBuilder().
			Context(ctx).
			Proxies(c.proxies).
			GET("https://you.com/api/user/getYouProState").
			Header("Cookie", c.cookie).
			DoWith(http.StatusOK)
		if err != nil {
			return nil, err
		}

		type state struct {
			Freemium      map[string]int
			Subscriptions []interface{}
		}

		var s state
		if err = common.ToObject(response, &s); err != nil {
			return nil, err
		}

		if s.Freemium["max_calls"] == s.Freemium["used_calls"] {
			return nil, errors.New("zero quota")
		}
	}

	var files []byte
	if previousMessages != "" {
		filename, err := upload(ctx, c.cookie, c.proxies, previousMessages)
		if err != nil {
			return nil, err
		}

		file := map[string]string{
			"user_filename": "paste.txt",
			"filename":      filename,
			"size":          strconv.Itoa(len(previousMessages)),
		}

		files, err = json.Marshal([]interface{}{file})
		if err != nil {
			return nil, err
		}
	}

	chatId := uuid.NewString()
	conversationTurnId := uuid.NewString()
	response, err := common.ClientBuilder().
		GET("https://you.com/api/streamingSearch").
		Context(ctx).
		Proxies(c.proxies).
		Query("q", url.QueryEscape(query)).
		Query("page", "1").
		Query("count", "10").
		Query("incognito", "true").
		Query("safeSearch", "off").
		Query("mkt", "zh-CN").
		Query("responseFilter", "TimeZone,Computation,RelatedSearches").
		Query("domain", "youchat").
		Query("use_personalization_extraction", "true").
		Query("chatId", chatId).
		Query("queryTraceId", chatId).
		Query("conversationTurnId", conversationTurnId).
		Query("pastChatLength", "0").
		Query("selectedChatMode", "custom").
		Query("userFiles", string(files)).
		Query("selectedAIModel", c.model).
		Query("traceId", fmt.Sprintf("%s|%s|%s", chatId, conversationTurnId, time.Now().Format(time.RFC3339)+"Z")).
		Query("chat", "[]").
		Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0").
		Header("Origin", "https://you.com").
		Header("Referer", "https://you.com/search").
		Header("Accept", "text/event-stream").
		Header("Priority", "u=1, i").
		Header("Cookie", fmt.Sprintf("ai_model=%s; %s", c.model, c.cookie)).
		DoWith(http.StatusOK)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go c.resolve(ctx, ch, response)
	return ch, nil
}

// 额度用完是否返回错误
func (c *Chat) LimitsWithE(limitsWithE bool) {
	c.limitsWithE = limitsWithE
}

// 附件上传
func upload(ctx context.Context, cookie, proxies, content string) (string, error) {
	response, err := common.ClientBuilder().
		Context(ctx).
		Proxies(proxies).
		GET("https://you.com/api/get_nonce").
		Header("Cookie", cookie).
		DoWith(http.StatusOK)
	if err != nil {
		return "", err
	}

	bio, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	file, _ := w.CreateFormFile("file", "paste.txt")
	file.Write([]byte(content))
	w.Close()

	response, err = common.ClientBuilder().
		Context(ctx).
		Proxies(proxies).
		POST("https://you.com/api/upload").
		Header("Cookie", cookie).
		Header("X-Upload-Nonce", string(bio)).
		Header("Content-Type", w.FormDataContentType()).
		Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0").
		Header("Origin", "https://you.com").
		Header("Referer", "https://you.com/search").
		Header("Priority", "u=1, i").
		SetBytes(b.Bytes()).
		DoWith(http.StatusOK)
	if err != nil {
		return "", err
	}

	var obj map[string]string
	if err = common.ToObject(response, &obj); err != nil {
		return "", err
	}

	if filename, ok := obj["filename"]; ok {
		return filename, nil
	}

	return "", errors.New("upload failed")
}

func (c *Chat) resolve(ctx context.Context, ch chan string, response *http.Response) {
	defer close(ch)

	scanner := bufio.NewScanner(response.Body)
	scanner.Split(func(data []byte, eof bool) (advance int, token []byte, err error) {
		if eof && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if eof {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	type chatToken struct {
		YouChatToken string `json:"youChatToken"`
	}

	// true 继续，false 结束
	eventHandler := func() bool {
		if !scanner.Scan() {
			return false
		}

		var event string
		data := scanner.Text()
		if data == "" {
			return true
		}

		if len(data) < 7 || data[:7] != "event: " {
			return true
		}
		event = data[7:]

		if event == "done" {
			return false
		}

		if !scanner.Scan() {
			return false
		}

		data = scanner.Text()
		if len(data) < 6 || data[:6] != "data: " {
			return true
		}
		data = data[6:]
		if event == "youChatModeLimits" {
			ch <- "limits: " + data
			return true
		}

		if event != "youChatToken" {
			return true
		}

		var token chatToken
		if err := json.Unmarshal([]byte(data), &token); err != nil {
			return true
		}

		if freeQuota(token.YouChatToken) {
			return true
		}

		ch <- token.YouChatToken
		return true
	}

	for {
		select {
		case <-ctx.Done():
			ch <- "error: context canceled"
			return
		default:
			if !eventHandler() {
				return
			}
		}
	}
}

func MergeMessages(messages []Message) string {
	if len(messages) == 0 {
		return ""
	}

	buffer := new(bytes.Buffer)
	lastRole := ""

	for _, message := range messages {
		if lastRole == "" || lastRole != message.Role {
			lastRole = message.Role
			buffer.WriteString(fmt.Sprintf("\n%s: %s", message.Role, message.Content))
			continue
		}
		buffer.WriteString(fmt.Sprintf("\n%s", message.Content))
	}

	return buffer.String()
}

func freeQuota(value string) bool {
	return strings.HasPrefix(value, "#### Please log in to access GPT-4 mode.") ||
		strings.HasPrefix(value, "#### You've hit your free quota for GPT-4 mode.")
}

func extCookies(cookies, model string) string {
	var result []string
	slice := strings.Split(cookies, "; ")
	for _, cookie := range slice {
		kv := strings.Split(cookie, "=")
		if len(kv) < 1 {
			continue
		}

		k := strings.TrimSpace(kv[0])
		v := strings.Join(kv[1:], "=")

		if strings.HasPrefix(k, "safesearch") {
			continue
		}

		if k == "ai_model" {
			result = append(result, k+"="+model)
			continue
		}

		result = append(result, k+"="+strings.TrimSpace(v))
	}

	return strings.Join(result, "; ")
}
