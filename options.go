package xnetwork

import (
	"github.com/xiecat/xhttp/xtls"
	"golang.org/x/time/rate"
)

var NetOptions *ClientOptions

type ClientOptions struct {
	FailRetries int                 `json:"fail_retries" yaml:"fail_retries" #:"请求失败的重试次数，0 则不重试"`
	MaxQPS      int                 `json:"max_qps" yaml:"max_qps" #:"每秒最大请求数"`
	ReadSize    int                 `json:"read_size" yaml:"read_size" #:"响应读取长度, 默认 2048"`
	DialTimeout int                 `json:"dial_timeout" yaml:"dial_timeout" #:"建立 tcp 连接的超时时间"`
	KeepAlive   int                 `json:"keep_alive" yaml:"keep_alive" #:"tcp keep_alive 时间"`
	ReadTimeout int                 `json:"read_timeout" yaml:"read_timeout" #:"响应读取超时时间"`
	TlsOptions  *xtls.ClientOptions `json:"tls" yaml:"tls" #:"tls 配置"`
	Debug       bool                `json:"net_debug" yaml:"net_debug" #:"是否启用 debug 模式, 开启 request trace"`
	Limiter     *rate.Limiter       `json:"-" yaml:"-"`
}

func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Debug:       false,
		FailRetries: 0, // 默认改为0，否则如果配置文件指定了0，会不生效。 "nil value" 的问题
		ReadSize:    2048,
		MaxQPS:      500,
		DialTimeout: 10,
		KeepAlive:   10,
		ReadTimeout: 5,
		TlsOptions:  xtls.DefaultClientOptions(),
		Limiter:     rate.NewLimiter(500, 1),
	}
}

func (o *ClientOptions) SetLimiter() *ClientOptions {
	o.Limiter = rate.NewLimiter(rate.Limit(o.MaxQPS), 1)
	return o
}

func GetNetOptions() *ClientOptions {
	if NetOptions != nil {
		return NetOptions
	} else {
		return DefaultClientOptions()
	}
}
