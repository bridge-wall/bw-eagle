package eagle

import (
	"context"
	"errors"
	"time"

	"github.com/bridge-wall/bw-eagle/middleware"
	"github.com/bridge-wall/bw-eagle/utils"
	"github.com/gin-gonic/gin"
)

type HttpClient struct {
	*middleware.HttpClient
}

func NewHttpClient(t time.Duration) *HttpClient {
	return &HttpClient{middleware.NewHttpClient(t)}
}


func (client *HttpClient) Get(ctx context.Context, url string, query map[string]string, headers map[string]string) (string, error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.HttpGet(c, url, query, headers)
		}
	}

	return "", errors.New("context is not gin.Context")
}


// Content-Type 是 application/json;charset=utf-8
func (client *HttpClient) Post(ctx context.Context, url string, query map[string]string, headers map[string]string, body interface{}) (string, error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.JsonPost(c, url, query, headers, body)
		}
	}

	return "", errors.New("context is not gin.Context")
}

// 并发请求

// Content-Type 是 application/json;charset=utf-8
func (client *HttpClient) MultiplePost(ctx context.Context, url string, params map[string]string, headers map[string]string, bodyList map[int]interface{}) (map[int]string, map[int]error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.MultipleJsonPost(c, url, params, headers, bodyList)
		}
	}

	err := make(map[int]error)
	err[0] = errors.New("context is not gin.Context")
	return nil, err
}

// 请求第三方服务
// Content-Type 是 application/x-www-form-urlencoded
// 无链路追踪
func (client *HttpClient) TFormGet(ctx context.Context, addr string, query map[string]string, headers map[string]string) (string, error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.ThirdFormGet(c, addr, query, headers)
		}
	}

	return "", errors.New("context is not gin.Context")
}

// 请求第三方服务
// Content-Type 是 application/x-www-form-urlencoded
// 无链路追踪
func (client *HttpClient) TFormPost(ctx context.Context, addr string, query map[string]string, headers map[string]string, body map[string]string) (string, error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.ThirdFormPost(c, addr, query, headers, body)
		}
	}

	return "", errors.New("context is not gin.Context")
}


// Content-Type 自定义
func (client *HttpClient) OriginPost(ctx context.Context, url string, query map[string]string, headers map[string]string, body string) (string, error) {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return client.HttpOriginPost(c, url, query, headers, body)
		}
	}

	return "", errors.New("context is not gin.Context")
}

// 生成请求体签名
func BodySign(key, body string) (string, error) {
	if key == "" || body == "" {
		return "", errors.New("param empty")
	}

	sign := utils.Md5(key + body)
	return sign, nil
}
