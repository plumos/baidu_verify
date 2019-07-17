package content_verify

import (
	"boss/base"
	"encoding/json"
	"github.com/astaxie/beego"
)

// ========== 5203图像审核接口
type VerifyBigImgStoreRequest struct {
	VerifyImg    string `json:"verifyImg"`    //图像大小4M以下，最长边不大于4096
	VerifyImgUrl string `json:"verifyImgUrl"` //图像大小4M以下，最长边不大于4096
}

type VerifyBigImgSend struct {
	Image  string   `json:"image"`  //图像数据，base64编码，不能与imgUrl并存
	ImgUrl string   `json:"imgUrl"` //图像Url，不能与image并存，不需要urlEncode
	Scenes []string `json:"scenes"` //指定本次调用的模型服务，以字符串数组表示。
	//scenesConf
}

type VerifyBigImgStoreReply struct {
}

type VerifyBigImgStore struct {
	base.BaseData
	Base  base.BaseRequest         `json:"base"`
	Param VerifyBigImgStoreRequest `json:"param"`
}

func (contents *VerifyBigImgStore) DoDecode() (object interface{}, err *base.Error) {
	base.LOG_INFO("context:" + string(contents.Input.BeegoCtx.Input.RequestBody))
	e := json.Unmarshal(contents.Input.BeegoCtx.Input.RequestBody, contents)
	if e != nil {
		err = base.NewReplayError(base.CODE_UNKNOWN_JSON_ERROR, base.CODE_UNKNOWN_JSON_ERROR_MESSAGE, nil)
		return nil, err
	}
	return nil, nil
}

func (contents *VerifyBigImgStore) DoProcess() (object interface{}, err *base.Error) {

	var url = beego.AppConfig.String("baidu_verify::img_url")
	url = url + GetSetBaiduToken()
	//ctx, cancelFunc := context.WithTimeout(context.Background(), (3 * time.Second))
	//defer cancelFunc()
	//header := map[string]string{"Content-Type": "application/json"}
	option := HttpOption{
		Retry: 1,
		//Timeout: 15000,
		URL: url,
		//	Header: header,
	}
	formDatas := VerifyBigImgSend{
		Image:  contents.Param.VerifyImg,
		ImgUrl: contents.Param.VerifyImgUrl,
		Scenes: []string{"politician", "antiporn", "terror", "disgust", "watermark"},
	}
	//politician：政治敏感识别 ,antiporn：色情识别 ,terror：暴恐识别, disgust:恶心图像识别, watermark:广告检测
	var result interface{}
	rawResponse, httpHeaders, er := JSONPostWithOption(option, formDatas, &result)
	base.LOG_INFO("[Post5203]", rawResponse, httpHeaders, er, result)
	contents.Reply.Biz = result
	return nil, nil

}
