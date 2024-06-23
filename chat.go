package you

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bincooo/emit.io"
	"github.com/gingfrederik/docx"
	_ "github.com/gingfrederik/docx"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

type Chat struct {
	cookie     string
	clearance  string
	model      string
	proxies    string
	limitWithE bool

	session   *emit.Session
	userAgent string
}

const (
	GPT_4       = "gpt_4"
	GPT_4o      = "gpt_4o"
	GPT_4_TURBO = "gpt_4_turbo"

	CLAUDE_2          = "claude_2"
	CLAUDE_3_HAIKU    = "claude_3_haiku"
	CLAUDE_3_SONNET   = "claude_3_sonnet"
	CLAUDE_3_5_SONNET = "claude_3_5_sonnet"
	CLAUDE_3_OPUS     = "claude_3_opus"
)

func New(cookie, model, proxies string) Chat {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0"
	return Chat{cookie, "", model, proxies, false, nil, userAgent}
}

func (c *Chat) Client(session *emit.Session) {
	c.session = session
}

func (c *Chat) CloudFlare(cookie, userAgent string) {
	c.clearance = cookie
	c.userAgent = userAgent
}

func (c *Chat) Reply(ctx context.Context, previousMessages []Message, query string, isf bool) (chan string, error) {
	if c.clearance == "" && cmdPort != "" {
		response, err := emit.ClientBuilder(c.session).
			Context(ctx).
			GET("http://127.0.0.1:" + cmdPort + "/clearance").
			DoS(http.StatusOK)
		if err != nil {
			return nil, err
		}

		obj, err := emit.ToMap(response)
		if err != nil {
			return nil, err
		}

		data := obj["data"].(map[string]interface{})
		c.clearance = data["cookie"].(string)
		c.userAgent = data["userAgent"].(string)
	}

	jar := extCookies(emit.MergeCookies(c.cookie, c.clearance), c.model)
	if c.limitWithE {
		count, err := c.State(ctx)
		if err != nil {
			return nil, err
		}
		if count <= 0 {
			return nil, errors.New("ZERO QUOTA")
		}
	}

	messages, err := mergeMessages(previousMessages, isf)
	if err != nil {
		return nil, err
	}

	var (
		userFiles = "_"
		files     = ""
		chatL     = strconv.Itoa(len(previousMessages))
	)

	if isf {
		filename, e := c.upload(ctx, c.proxies, jar, messages)
		if e != nil {
			return nil, e
		}
		userFiles = "userFiles"
		files = fmt.Sprintf(`[{"user_filename":"messages.txt","filename":"%s","size":"%d"}]`, filename, len(messages))
		chatL = "0"
	}

	chatId := uuid.NewString()
	conversationTurnId := uuid.NewString()
	t := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")
	response, err := emit.ClientBuilder(c.session).
		GET("https://you.com/api/streamingSearch").
		Context(ctx).
		Proxies(c.proxies).
		CookieJar(jar).
		Query("q", url.QueryEscape(query)).
		Query("page", "1").
		Query("count", "10").
		Query("safeSearch", "Off").
		Query("mkt", "zh-HK").
		Query("domain", "youchat").
		Query("use_personalization_extraction", "false").
		Query("queryTraceId", chatId).
		Query("chatId", chatId).
		Query("conversationTurnId", conversationTurnId).
		Query("selectedChatMode", "custom").
		Query(userFiles, url.QueryEscape(files)).
		Query("selectedAiModel", c.model).
		Query("traceId", fmt.Sprintf("%s|%s|%s", chatId, conversationTurnId, t)).
		Query("incognito", "true").
		//Query("responseFilter", "WebPages,TimeZone,Computation,RelatedSearches").
		Query("pastChatLength", chatL).
		Query("chat", url.QueryEscape(messages)).
		Header("User-Agent", c.userAgent).
		Header("Host", "you.com").
		Header("Origin", "https://you.com").
		Header("Referer", "https://you.com/search?fromSearchBar=true&tbm=youchat&chatMode=custom").
		Header("Accept-Language", "en-US,en;q=0.9").
		Header("Accept", "text/event-stream").
		DoS(http.StatusOK)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)
	go c.resolve(ctx, ch, response)
	return ch, nil
}

func (c *Chat) State(ctx context.Context) (int, error) {
	response, err := emit.ClientBuilder(c.session).
		Context(ctx).
		Proxies(c.proxies).
		GET("https://you.com/api/user/getYouProState").
		Header("Cookie", emit.MergeCookies(c.cookie, c.clearance)).
		Header("User-Agent", c.userAgent).
		Header("Accept-Language", "en-US,en;q=0.9").
		Header("Referer", "https://you.com/").
		Header("Origin", "https://you.com").
		DoS(http.StatusOK)
	if err != nil {
		return -1, err
	}

	type state struct {
		Freemium      map[string]int
		Subscriptions []interface{}
	}

	var s state
	if err = emit.ToObject(response, &s); err != nil {
		return -1, err
	}

	logrus.Infof("used: %d/%d", s.Freemium["used_calls"], s.Freemium["max_calls"])
	if s.Freemium["max_calls"] == s.Freemium["used_calls"] {
		return 0, nil
	}

	return s.Freemium["max_calls"] - s.Freemium["used_calls"], nil
}

// 额度用完是否返回错误
func (c *Chat) LimitWithE(limitWithE bool) {
	c.limitWithE = limitWithE
}

// 附件上传
func (c *Chat) upload(ctx context.Context, proxies string, jar http.CookieJar, content string) (string, error) {
	response, err := emit.ClientBuilder(c.session).
		Context(ctx).
		Proxies(proxies).
		GET("https://you.com/api/get_nonce").
		CookieJar(jar).
		Header("Accept", "application/json, text/plain, */*").
		Header("Accept-Language", "en-US,en;q=0.9").
		Header("Referer", "https://you.com/?chatMode=custom").
		Header("Origin", "https://you.com").
		Header("User-Agent", c.userAgent).
		DoS(http.StatusOK)
	if err != nil {
		return "", err
	}

	uploadNonce := emit.TextResponse(response)

	doc := docx.NewFile()
	para := doc.AddParagraph()
	para.AddText(content)

	var buffer bytes.Buffer

	//h := make(textproto.MIMEHeader)
	//h.Set("Content-Disposition", `form-data; name="file"; filename="messages.docx"`)
	//h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	//h.Set("Content-Type", "text/plain")
	//fw, _ := w.CreatePart(h)
	//err = doc.Write(fw)
	//if err != nil {
	//	return "", err
	//}

	w := multipart.NewWriter(&buffer)
	fw, _ := w.CreateFormFile("file", "messages.txt")
	_, err = io.Copy(fw, strings.NewReader(content))
	if err != nil {
		return "", err
	}
	_ = w.Close()

	response, err = emit.ClientBuilder(c.session).
		Context(ctx).
		Proxies(proxies).
		CookieJar(jar).
		POST("https://you.com/api/upload").
		Header("X-Upload-Nonce", uploadNonce).
		Header("Content-Type", w.FormDataContentType()).
		Header("Origin", "https://you.com").
		Header("Accept-Language", "en-US,en;q=0.9").
		Header("Host", "you.com").
		Header("Accept-Encoding", "br").
		Header("Referer", "https://you.com/?chatMode=custom").
		Header("Origin", "https://you.com").
		Header("Accept", "multipart/form-data").
		Header("User-Agent", c.userAgent).
		Buffer(&buffer).
		DoS(http.StatusOK)
	if err != nil {
		return "", err
	}

	var obj map[string]string
	if err = emit.ToObject(response, &obj); err != nil {
		return "", err
	}

	if filename, ok := obj["filename"]; ok {
		response, err = emit.ClientBuilder(c.session).
			Context(ctx).
			Proxies(proxies).
			CookieJar(jar).
			POST("https://you.com/api/instrumentation").
			JHeader().
			Header("Origin", "https://you.com").
			Header("Accept-Language", "en-US,en;q=0.9").
			Header("Host", "you.com").
			Header("Accept-Encoding", "br").
			Header("Referer", "https://you.com/?chatMode=custom").
			Header("Origin", "https://you.com").
			Header("Accept", "application/json, text/plain, */*").
			Header("User-Agent", c.userAgent).
			Bytes([]byte(`{"metricName":"file_upload_client_info_file_drop","documentVisibilityState":"visible","metricType":"info","value":1}`)).
			DoS(http.StatusOK)
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

		logrus.Trace("--------- ORIGINAL MESSAGE ---------")
		logrus.Trace(data)

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
		logrus.Trace(data)
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

		if quotaEmpty(token.YouChatToken) {
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

func mergeMessages(messages []Message, files bool) (string, error) {
	if len(messages) == 0 {
		return "[]", nil
	}

	if files {
		var buffer bytes.Buffer
		messageL := len(messages)
		for pos, message := range messages {
			buffer.WriteString(fmt.Sprintf("%s\n\n%s", message.Question, message.Answer))
			if pos < messageL-1 {
				buffer.WriteString("\n\n")
			}
		}
		return buffer.String(), nil
	}

	messageBytes, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}

	return string(messageBytes), nil
}

func quotaEmpty(value string) bool {
	return strings.HasPrefix(value, "#### Please log in to access GPT-4 mode.") ||
		strings.HasPrefix(value, "#### You've hit your free quota for GPT-4 mode.")
}

func extCookies(cookies, model string) (jar http.CookieJar) {
	jar, _ = cookiejar.New(nil)
	u, _ := url.Parse("https://you.com")

	slice := strings.Split(cookies, "; ")
	for _, cookie := range slice {
		kv := strings.Split(cookie, "=")
		if len(kv) < 1 {
			continue
		}

		k := strings.TrimSpace(kv[0])
		v := strings.Join(kv[1:], "=")

		if strings.HasPrefix(k, "safesearch") {
			jar.SetCookies(u, []*http.Cookie{{Name: k, Value: "Off"}})
			continue
		}

		if k == "you_subscription" {
			jar.SetCookies(u, []*http.Cookie{{Name: k, Value: "freemium"}})
			continue
		}

		if k == "ai_model" {
			jar.SetCookies(u, []*http.Cookie{{Name: k, Value: model}})
			continue
		}

		jar.SetCookies(u, []*http.Cookie{{Name: k, Value: strings.TrimSpace(v)}})
	}

	//
	jar.SetCookies(u, []*http.Cookie{{Name: "has_seen_agent_uploads_modal", Value: "true"}})
	return
}
