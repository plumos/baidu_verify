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

// ========== 5202头像审核接口
type VerifyImgStoreRequest struct {
	VerifyImg    string `json:"verifyImg"`    //图像大小4M以下，最长边不大于4096
	VerifyImgUrl string `json:"verifyImgUrl"` //图像大小4M以下，最长边不大于4096
}

type VerifyImgStoreReply struct {
	//	LogId  int64                   `json:"log_id"` //正确调用生成的唯一标识码，用于问题定位
	//	Result []VerifyImgContentReply `json:"result"`
}

type VerifyImgContentReply struct {
	//Res_msg  []int          `json:"res_msg"`  //未校验通过的项,[]描述的是未校验通过的规则，具体参考 result中的res_msg业务错误码定义
	//Res_code int64          `json:"res_code"` //业务校验结果 0：校验通过，1：校验不通过
	//Data     VerifyImgClass `json:"data"`     //识别详细结果数据
}

type VerifyImgClass struct {
	Antiporn   VerifyImgData `json:"antiporn"`   //色情识别返回结果
	Face       VerifyImgData `json:"face"`       //人脸检测服务返回结果
	Terror     VerifyImgData `json:"terror"`     //暴恐识别返回结果
	Public     VerifyImgData `json:"public"`     //公众人物服务返回结果
	Ocr        VerifyImgData `json:"ocr"`        //文字识别服务返回结果
	Politician VerifyImgData `json:"politician"` //政治敏感识别返回结果
	Quality    VerifyImgData `json:"quality"`    //图像质量返回结果
}

type VerifyImgData struct {
	Result      []ImgResult `json:"result"`
	LogId       int64       `json:"log_id"`      //请求标识码，随机数，唯一
	Error_msg   string      `json:"error_msg"`   //错误信息，参考错误码表错误码说明。只在异常响应中出现
	Error_code  int64       `json:"error_code"`  //错误码，参考错误码表错误码说明。只在异常响应中出现
	Conclusion  string      `json:"conclusion"`  //本张图片最终鉴定的结果，分为“色情”，“性感”，“正常”三种。
	Result_fine []ImgResult `json:"result_fine"` //对应标签的置信度得分，越高可信度越高
}

type ImgResult struct {
	Class_name  string  `json:"class_name"`  // 	分类结果名称
	Probability float64 `json:"probability"` //错误信息，参考错误码表错误码说明。只在异常响应中出现
}

type VerifyImgContent struct {
	Score float64  `json:"score"` //违禁检测分，范围0~1，数值从低到高代表风险程度的高低
	Hit   []string `json:"hit"`   //违禁类型对应命中的违禁词集合，可能为空
	Label int      `json:"label"` //请求中的违禁类型:1.暴恐违禁,2.文本色情,3.政治敏感,4.恶意推广,5.低俗辱骂,6.低质灌水
}

type VerifyImgStore struct {
	base.BaseData
	Base  base.BaseRequest      `json:"base"`
	Param VerifyImgStoreRequest `json:"param"`
}

func (contents *VerifyImgStore) DoDecode() (object interface{}, err *base.Error) {
	base.LOG_INFO("context:" + string(contents.Input.BeegoCtx.Input.RequestBody))
	e := json.Unmarshal(contents.Input.BeegoCtx.Input.RequestBody, contents)
	if e != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_JSON_ERROR, base.CODE_UNKNOWN_JSON_ERROR_MESSAGE, nil)
		return nil, err
	}
	return nil, nil
}

func (contents *VerifyImgStore) DoProcess() (object interface{}, err *base.Error) {

	client := &http.Client{}
	var reply interface{}
	var er error
	var req *http.Request
	var url = beego.AppConfig.String("baidu_verify::logo_url")
	url = url + GetSetBaiduToken()
	if !base.IsEmpty(contents.Param.VerifyImg) {
		s, head_err := base.Base64Decode(contents.Param.VerifyImg, 1)
		if head_err != nil {
			base.LOG_ERROR("head_err : ", head_err)
			err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
			return nil, err
		}
		req, er = http.NewRequest("POST", url, strings.NewReader("images="+s))
	} else {
		s := base.EncodeURL(contents.Param.VerifyImgUrl)
		req, er = http.NewRequest("POST", url, strings.NewReader("imgUrls="+s))
	}
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
	//err = base.NewReplayError(base.CODE_VERIFY_OVER_LIMIT, base.CODE_VERIFY_OVER_LIMIT_MESSAGE, nil)
	//return nil, err
	return nil, nil

}
