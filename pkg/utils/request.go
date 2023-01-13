package utils

import (
	"Yi/pkg/logging"
	"crypto/tls"
	"go.uber.org/ratelimit"
	"net/http"
	"net/url"
)

/**
  @author: yhy
  @since: 2023/1/12
  @desc: //TODO
**/

type Session struct {
	// Client is the current http client
	Client *http.Client
	// Rate limit instance
	RateLimiter ratelimit.Limiter // 每秒请求速率限制
}

func NewSession(proxy string) *Session {
	Transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		DisableKeepAlives:   true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Add proxy
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		if isSupportedProtocol(proxyURL.Scheme) {
			Transport.Proxy = http.ProxyURL(proxyURL)
		} else {
			logging.Logger.Warnln("Unsupported proxy protocol: %s", proxyURL.Scheme)
		}
	}

	client := &http.Client{
		Transport: Transport,
	}
	session := &Session{
		Client: client,
	}

	// github api 访问加上 token 的访问速率为每小时 5000 次，平均下来每秒一次多，这里限制为每秒访问一次 github
	session.RateLimiter = ratelimit.New(1)

	return session

}
