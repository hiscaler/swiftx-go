package swiftx

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/swiftx-go/config"
	"github.com/hiscaler/swiftx-go/entity"
	"github.com/hiscaler/swiftx-go/response"
)

const (
	Version   = "0.0.1"
	userAgent = "SwiftX Express API Client-Golang/" + Version + " (https://github.com/hiscaler/swiftx-go)"
)

const (
	ProdBaseUrl = "https://prod.open.swiftx-express.com/api/v2/openapi"
	TestBaseUrl = "https://test.open.swiftx-express.com/api/v2/openapi"
)

type Client struct {
	config     *config.Config // 配置
	logger     *slog.Logger   // Logger
	httpClient *resty.Client  // Resty Client
	Services   services       // API Services
}

type signature struct {
	timestamp     int64
	nonce         string
	contentSHA256 string
	signature     string
}

// buildSignature 构建签名
// 签名格式：{app_key}\n{timestamp}\n{nonce}\n{content_sha256}\n{http_method}\n{path}\n{query_string}
func buildSignature(appKey, appSecret, httpMethod, apiPath, queryString string, requestBody any) signature {
	timestamp := time.Now().Unix()
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		panic("failed to generate random nonce: " + err.Error())
	}
	nonce := hex.EncodeToString(nonceBytes)

	// 计算 requestBody 的 SHA256 哈希值
	var bodyBytes []byte
	if requestBody != nil {
		bodyBytes, _ = json.Marshal(requestBody)
	}
	hasher := sha256.New()
	hasher.Write(bodyBytes)
	contentSHA256 := hex.EncodeToString(hasher.Sum(nil))

	// 构建待签名字符串
	stringToSign := fmt.Sprintf("%s\n%d\n%s\n%s\n%s\n%s\n%s",
		appKey,
		timestamp,
		nonce,
		contentSHA256,
		httpMethod,
		apiPath,
		queryString,
	)

	// 计算签名
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(stringToSign))
	sign := hex.EncodeToString(h.Sum(nil))

	return signature{
		timestamp:     timestamp,
		nonce:         nonce,
		contentSHA256: contentSHA256,
		signature:     sign,
	}
}

func NewClient(cfg config.Config) *Client {
	l := createLogger()
	debug := cfg.Debug
	if cfg.Logger != nil {
		l.l = cfg.Logger
	}

	swiftxClient := &Client{
		config: &cfg,
		logger: l.l,
	}
	baseUrl := ProdBaseUrl
	if cfg.Env != entity.Prod {
		baseUrl = TestBaseUrl
	}
	httpClient := resty.New().
		SetDebug(debug).
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
			"X-App-Key":    cfg.AppKey,
		}).
		SetTimeout(time.Duration(cfg.Timeout) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			u, err := url.Parse(request.URL)
			if err != nil {
				return err
			}
			sign := buildSignature(cfg.AppKey, cfg.AppSecret, request.Method, "/api/v2/openapi"+u.Path, request.QueryParam.Encode(), request.Body)
			request.SetHeaders(map[string]string{
				"X-Timestamp":      strconv.Itoa(int(sign.timestamp)),
				"X-Nonce":          sign.nonce,
				"X-Content-SHA256": sign.contentSHA256,
				"X-Signature":      sign.signature,
			})
			return nil
		}).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response != nil && response.StatusCode() == http.StatusTooManyRequests {
				return true
			}
			return false
		})
	swiftxClient.httpClient = httpClient
	xService := service{
		config:     &cfg,
		logger:     l.l,
		httpClient: swiftxClient.httpClient,
	}
	swiftxClient.Services = services{
		Order: (orderService)(xService),
		Ping:  (pingService)(xService),
	}
	return swiftxClient
}

// errorWrap 错误包装
func errorWrap(code int, message string) error {
	if code == http.StatusOK {
		return nil
	}

	switch code {
	case 401:
		message = "身份验证失败（无效签名）"
	case 403:
		message = "授权失败（权限不足）"
	case 429:
		message = "超出速率限制"
	case 500:
		message = "服务器错误，请联系 SwiftX Express 客服人员"
	default:
		message = strings.TrimSpace(message)
		if message == "" {
			message = "未知错误"
		}
	}
	return errors.New(message)
}

func invalidInput(e error) error {
	var errs validation.Errors
	if !errors.As(e, &errs) {
		return e
	}

	if len(errs) == 0 {
		return nil
	}

	fields := make([]string, 0)
	messages := make([]string, 0)
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		e1 := errs[field]
		if e1 == nil {
			continue
		}

		var errObj validation.ErrorObject
		if errors.As(e1, &errObj) {
			e1 = errObj
		} else {
			var errs1 validation.Errors
			if errors.As(e1, &errs1) {
				e1 = invalidInput(errs1)
				if e1 == nil {
					continue
				}
			}
		}

		messages = append(messages, e1.Error())
	}
	return errors.New(strings.Join(messages, "; "))
}

func recheckError(resp *resty.Response, e error) error {
	if e != nil {
		if errors.Is(e, http.ErrHandlerTimeout) {
			return errorWrap(http.StatusRequestTimeout, e.Error())
		}
		return e
	}

	if resp.IsSuccess() {
		return nil
	}

	if resp.IsError() {
		return errorWrap(resp.StatusCode(), "")
	}

	var normalResponse response.NormalResponse
	err := json.Unmarshal(resp.Body(), &normalResponse)
	if err != nil {
		return err
	}
	if normalResponse.Result.Success {
		return nil
	}
	return errors.New(normalResponse.Result.Message)
}
