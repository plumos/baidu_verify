package content_verify

import (
	"boss/base"
	"boss/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"time"
)

type BaiduTokenReply struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` //过期时间一个月
}

type HttpOption struct {
	Retry   int
	Timeout int // 单位,毫秒
	URL     string
	Header  map[string]string
}

type HttpResponse struct {
	StatusCode  int    // 状态码
	RawResponse []byte // 原始响应
}

type HttpHeaders struct {
	RequestHeader  map[string][]string // 请求header
	ResponseHeader map[string][]string // 响应接口
}

func init() {
	GetSetBaiduToken()
}

func GetSetBaiduToken() string {
	var key = VERIFY_BAIDU_KEY
	val := models.RedisGet(key)
	if val == nil {
		accessToken := GetBaiduToken()
		err := models.RedisPut(key, accessToken, 2590000*time.Second) //比百度默认时间短
		if err != nil {
			base.LOG_ERROR("accessToken RedisPut err : ", err)
		}
		return accessToken
	}
	return base.B2S(val.([]uint8))
}

func GetBaiduToken() string {

	client := &http.Client{}
	var config BaiduTokenReply
	var beforeUrl = beego.AppConfig.String("baidu_verify::token_url")
	req, er := http.NewRequest("POST", beforeUrl, nil)
	if er != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, er := client.Do(req)
	defer resp.Body.Close()

	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		base.LOG_ERROR("handle error : ", er) // handle error
	}
	e := json.Unmarshal(body, &config)
	if e != nil {
		base.LOG_ERROR("json error : ", er) // handle error
	}
	return config.AccessToken
}

func StartPostContext(ctx context.Context, option *HttpOption, body *bytes.Buffer, dest interface{}) (
	rawHttpResponse HttpResponse, httpHeaders HttpHeaders, err error) {
	retry := option.Retry
	if retry <= 0 {
		retry = 1
	}
	timeout := time.Duration(time.Duration(option.Timeout) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	var request *http.Request
	if request, err = http.NewRequest("POST", option.URL, body); err != nil {
		return
	}
	request = request.WithContext(ctx)
	// 处理头部
	for k, v := range option.Header {
		request.Header.Set(k, v)
	}
	httpHeaders = HttpHeaders{
		RequestHeader: request.Header,
	}

	for i := 0; i < retry; i++ {
		var rawResponse *http.Response
		rawResponse, err = client.Do(request)

		if rawResponse != nil {
			httpHeaders.ResponseHeader = rawResponse.Header
		}

		if err != nil {
			continue
		}

		defer rawResponse.Body.Close()
		var respData []byte
		respData, err = ioutil.ReadAll(rawResponse.Body)

		if err != nil {
			continue
		}

		rawHttpResponse = HttpResponse{
			StatusCode:  rawResponse.StatusCode,
			RawResponse: respData,
		}
		if rawResponse.StatusCode != 200 {
			base.LOG_ERROR("Post Return Http Status:", rawResponse.Status)
			err = errors.New("Post Return Http Status :" + rawResponse.Status)
			continue
		}
		base.LOG_INFO("[Post] URL: %s, Request: %s, Response: %s", option.URL, body, respData)

		err = json.Unmarshal(respData, dest)
		if err != nil {
			base.LOG_ERROR("Post Return Http Status:", err.Error())
			continue
		}
		if err == nil {
			break
		}
	}
	return
}

// JSONPost 用json向一个URL发送请求，并解析返回值
func JSONPostWithOption(option HttpOption, data interface{}, dest interface{}) (rawHttpResponse HttpResponse, httpHeaders HttpHeaders, err error) {
	retry := option.Retry
	if retry <= 0 {
		retry = 1
	}

	var body []byte
	if body, err = json.Marshal(data); err != nil {
		return
	}
	timeout := time.Duration(time.Duration(option.Timeout) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	for i := 0; i < retry; i++ {
		var request *http.Request
		request, err = http.NewRequest("POST", option.URL, bytes.NewBuffer(body))
		if err != nil {
			continue
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		var rawResponse *http.Response
		rawResponse, err = client.Do(request)

		// 添加记录headers
		httpHeaders = HttpHeaders{
			RequestHeader: request.Header,
		}
		if rawResponse != nil {
			httpHeaders.ResponseHeader = rawResponse.Header
		}

		if err != nil {
			continue
		}

		defer rawResponse.Body.Close()
		var respData []byte
		respData, err = ioutil.ReadAll(rawResponse.Body)
		if err != nil {
			continue
		}
		rawHttpResponse = HttpResponse{
			StatusCode:  rawResponse.StatusCode,
			RawResponse: respData,
		}
		if rawResponse.StatusCode != 200 {
			err = errors.New("JSONPost Return Http Status :" + rawResponse.Status)
			continue
		}

		base.LOG_INFO("[Post] URL: %s, Request: %s, Response: %s", option.URL, body, respData)

		err = json.Unmarshal(respData, dest)
		if err != nil {
			base.LOG_ERROR("Post Return Http Status:", err.Error())
			continue
		}
		if err == nil {
			break
		}
	}
	return
}
