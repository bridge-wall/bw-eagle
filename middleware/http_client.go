package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpClient struct {
	timeout time.Duration
}

func NewHttpClient(t time.Duration) *HttpClient {
	return &HttpClient{
		timeout: t,
	}
}

func (cli *HttpClient) addQuery(req *http.Request, params map[string]string) {
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func (cli *HttpClient) setHeaders(req *http.Request, headers map[string]string) {
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}
func (cli *HttpClient) setHeader(req *http.Request, key, value string) {
	req.Header.Set(key, value)
}


func (cli *HttpClient) HttpGet(c *gin.Context, addr string, query map[string]string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return "", err
	}

	cli.addQuery(req, query)
	cli.setHeaders(req, headers)


	client := &http.Client{
		Timeout: cli.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}

func (cli *HttpClient) JsonPost(c *gin.Context, addr string, query map[string]string, headers map[string]string, body interface{}) (string, error) {
	var err error
	var bodyJson []byte
	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			return "", err
		}
	}

	req, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(bodyJson)))
	if err != nil {
		return "", err
	}

	cli.addQuery(req, query)
	cli.setHeaders(req, headers)
	cli.setHeader(req, "Content-Type", "application/json;charset=utf-8")

	client := &http.Client{
		Timeout: cli.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}

// 并发请求
func (cli *HttpClient) MultipleJsonPost(c *gin.Context, addr string, params map[string]string, headers map[string]string, bodyList map[int]interface{}) (map[int]string, map[int]error) {

	// 非线程安全，需要加锁
	var respLock sync.Mutex
	var errLock sync.Mutex
	respList := make(map[int]string)
	errList := make(map[int]error)

	var wg sync.WaitGroup
	for k, body := range bodyList {
		wg.Add(1)
		go func(k int, body interface{}) error {
			defer func() {
				if err := recover(); err != nil {
					LogError(nil, "[RecoveryPanic] [HttpClient] "+err.(string))
				}
			}()

			defer wg.Done()

			var bodyJson []byte
			if body != nil {
				var err error
				bodyJson, err = json.Marshal(body)
				if err != nil {
					errLock.Lock()
					errList[k] = err
					errLock.Unlock()
					return err
				}
			}

			req, err := http.NewRequest("POST", addr, bytes.NewBuffer(bodyJson))
			if err != nil {
				errLock.Lock()
				errList[k] = err
				errLock.Unlock()
				return err
			}

			cli.addQuery(req, params)
			cli.setHeaders(req, headers)
			cli.setHeader(req, "Content-Type", "application/json;charset=utf-8")

			client := &http.Client{
				Timeout: cli.timeout,
			}
			resp, err := client.Do(req)
			if err != nil {
				errLock.Lock()
				errList[k] = err
				errLock.Unlock()
				return err
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errLock.Lock()
				errList[k] = err
				errLock.Unlock()
				return err
			}

			if resp.StatusCode != http.StatusOK {
				err := errors.New(string(respBody))
				errLock.Lock()
				errList[k] = err
				errLock.Unlock()
				return err
			}

			respLock.Lock()
			respList[k] = string(respBody)
			respLock.Unlock()

			return nil

		}(k, body)
	}

	wg.Wait()

	return respList, errList
}

func (cli *HttpClient) ThirdFormGet(c *gin.Context, addr string, query map[string]string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return "", err
	}

	cli.addQuery(req, query)
	cli.setHeaders(req, headers)
	cli.setHeader(req, "Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: cli.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}

func (cli *HttpClient) ThirdFormPost(c *gin.Context, addr string, query map[string]string, headers map[string]string, body map[string]string) (string, error) {
	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", addr, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	cli.addQuery(req, query)
	cli.setHeaders(req, headers)
	cli.setHeader(req, "Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: cli.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}

func (cli *HttpClient) HttpOriginPost(c *gin.Context, addr string, query map[string]string, headers map[string]string, body string) (string, error) {
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	cli.addQuery(req, query)
	cli.setHeaders(req, headers)

	client := &http.Client{
		Timeout: cli.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}
