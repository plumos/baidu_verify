package content_verify

import (
	"boss/base"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
)

// ========== 5201评论审核接口
type ContentVerifyStoreRequest struct {
	VerifyText string `json:"verifyText"`
}

type ContentVerifyStoreReply struct {
	//	LogId  int64             `json:"log_id"` //正确调用生成的唯一标识码，用于问题定位
	//	Result BaiduVerifynReply `json:"result"`
}

type BaiduVerifynReply struct {
	Spam   int             `json:"spam"`   //请求中是否包含违禁，0表示非违禁，1表示违禁，2表示建议人工复审
	Reject []VerifyContent `json:"reject"` //审核未通过的类别列表与详情
	Review []VerifyContent `json:"review"` //待人工复审的类别列表与详情
	Pass   []VerifyContent `json:"pass"`   //审核通过的类别列表与详情
}

type VerifyContent struct {
	Score float64  `json:"score"` //违禁检测分，范围0~1，数值从低到高代表风险程度的高低
	Hit   []string `json:"hit"`   //违禁类型对应命中的违禁词集合，可能为空
	Label int      `json:"label"` //请求中的违禁类型:1.暴恐违禁,2.文本色情,3.政治敏感,4.恶意推广,5.低俗辱骂,6.低质灌水
}

type ContentVerifyStore struct {
	base.BaseData
	Base  base.BaseRequest          `json:"base"`
	Param ContentVerifyStoreRequest `json:"param"`
}

func (contents *ContentVerifyStore) DoDecode() (object interface{}, err *base.Error) {
	base.LOG_INFO("context:" + string(contents.Input.BeegoCtx.Input.RequestBody))
	e := json.Unmarshal(contents.Input.BeegoCtx.Input.RequestBody, contents)
	if e != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_JSON_ERROR, base.CODE_UNKNOWN_JSON_ERROR_MESSAGE, nil)
		return nil, err
	}
	return nil, nil
}

func (contents *ContentVerifyStore) DoProcess() (object interface{}, err *base.Error) {
	client := &http.Client{}
	var reply interface{}
	var url = beego.AppConfig.String("baidu_verify::text_url")
	url = url + GetSetBaiduToken()
	req, er := http.NewRequest("POST", url, strings.NewReader("content="+contents.Param.VerifyText))
	if er != nil {
		err = base.NewReplayError(base.CODE_PARAM_ERROR, base.CODE_PARAM_ERROR_MESSAGE, nil)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, er := client.Do(req)
	defer resp.Body.Close()

	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
		return nil, err
	}
	fmt.Println(string(body))
	e := json.Unmarshal(body, &reply)
	if e != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
		return nil, err
	}
	contents.Reply.Biz = reply
	return nil, nil

}
