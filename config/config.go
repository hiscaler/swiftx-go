package config

type Config struct {
	Debug       bool   `json:"debug"`        // 是否启用调试模式
	Env         string `json:"env"`          // 环境
	Timeout     int    `json:"timeout"`      // HTTP 超时设定（单位：秒）
	AppKey      string `json:"account"`      // 应用程序的唯一标识符
	AppSecret   string `json:"password"`     // 密钥
	CallbackUrl string `json:"callback_url"` // 回调地址
}
