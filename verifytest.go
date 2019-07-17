package content_verify

import (
	"boss/base"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ========== 5201评论审核接口
type ContentsCenterStoreRequest struct {
	Typ         string 	`json:"type"`
	Cid         string 	`json:"cid"`
	CommentText string 	`json:"commentText"`
	CommentType string 	`json:"commentType"`
	ReplyId     string 	`json:"replyId"`
}

type ContentsCenterStoreReply struct {
	CommentId string 	`json:"commentId"`
	CommentTyp string 	`json:"commentType"`
}

type ContentsCenterStore struct {
	base.BaseData
	Base  base.BaseRequest           `json:"base"`
	Param ContentsCenterStoreRequest `json:"param"`
}

func (contents *ContentsCenterStore) DoDecode() (object interface{}, err *base.Error) {
	base.LOG_INFO("context:" + string(contents.Input.BeegoCtx.Input.RequestBody))
	e := json.Unmarshal(contents.Input.BeegoCtx.Input.RequestBody, contents)
	if e != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_JSON_ERROR, base.CODE_UNKNOWN_JSON_ERROR_MESSAGE, nil)
		return nil, err
	}
	return nil, nil
}

func (contents *ContentsCenterStore) DoProcess() (object interface{}, err *base.Error) {

	return  nil,nil
}

func httpDos() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}
	xBceDate := time.Now().Format("2006-01-02T08:23:49Z")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("host", "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=6v8YySztsiTooqzLlF6zwVIm&client_secret=5mXTDWr8k2fTGUTZnwYiK0SEg07Guf9x")
	req.Header.Set("x-bce-date", xBceDate)
	req.Header.Set("authorization", time.Now().Format("2006-01-02T08:23:49Z"))
	//CanonicalHeaders：对HTTP请求中的Header部分进行选择性编码的结果。
	//若URL为https://bos.cn-n1.baidubce.com/example/测试，则其URL Path为/example/测试，将之规范化得到CanonicalURI ＝ /example/%E6%B5%8B%E8%AF%95。
	//CanonicalQueryString：对于URL中的Query String（Query String即URL中“？”后面的“key1 = valve1 & key2 = valve2 ”字符串）进行编码后的结果。
	//CanonicalURI = UriEncodeExceptSlash(Path)
	//CanonicalRequest = HTTP Method + "\n" + CanonicalURI + "\n" + CanonicalQueryString + "\n" + CanonicalHeaders
	//SigningKey = HMAC-SHA256-HEX(sk, authStringPrefix)
	//Signature = HMAC-SHA256-HEX(SigningKey, CanonicalRequest)

	//CanonicalRequest := "POST" + "\n" + "/rest/2.0/solution/v1/img_censor/v2/user_defined" + "\n"
	//authorization := "bce-auth-v1/6v8YySztsiTooqzLlF6zwVIm/" + xBceDate + "/1800/host;x-bce-date/"
					//bce-auth-v1/46bd9968a6194b4bbdf0341f2286ccce/2015-03-24T13:02:00Z/1800/host;x-bce-date/994014d96b0eb26578e039fa053a4f9003425da4bfedf33f4790882fb4c54903
	//bce-auth-v1/{accessKeyId}/{timestamp}/{expirationPeriodInSeconds}/{signedHeaders}/{signature}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}