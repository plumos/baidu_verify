package content_verify

//
//import (
//	"boss/base"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"strings"
//)
//
//// ========== 5203图像审核接口
//type VerifyBigImgStoreRequest struct {
//	VerifyImg    string `json:"verifyImg"`    //图像大小4M以下，最长边不大于4096
//	VerifyImgUrl string `json:"verifyImgUrl"` //图像大小4M以下，最长边不大于4096
//}
//
//type VerifyBigImgStoreReply struct {
//}
//
//type VerifyBigImgStore struct {
//	base.BaseData
//	Base  base.BaseRequest         `json:"base"`
//	Param VerifyBigImgStoreRequest `json:"param"`
//}
//
//func (contents *VerifyBigImgStore) DoDecode() (object interface{}, err *base.Error) {
//	base.LOG_INFO("context:" + string(contents.Input.BeegoCtx.Input.RequestBody))
//	e := json.Unmarshal(contents.Input.BeegoCtx.Input.RequestBody, contents)
//	if e != nil {
//		err = base.NewReplayError(base.CODE_UNKNOWN_JSON_ERROR, base.CODE_UNKNOWN_JSON_ERROR_MESSAGE, nil)
//		return nil, err
//	}
//	return nil, nil
//}
//
//func (contents *VerifyBigImgStore) DoProcess() (object interface{}, err *base.Error) {
//
//	client := &http.Client{}
//	var reply interface{}
//	var er error
//	var req *http.Request
//	var url = "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/user_defined?access_token="
//	url = url + GetBaiduToken()
//	if !base.IsEmpty(contents.Param.VerifyImg) {
//		s, head_err := base.Base64Decode(contents.Param.VerifyImg, 1)
//		if head_err != nil {
//			// decode错误
//			base.LOG_INFO("head_err : ", head_err)
//			err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
//			return nil, err
//		}
//		req, er = http.NewRequest("POST", url, strings.NewReader("image="+s))
//	} else {
//		req, er = http.NewRequest("POST", url, strings.NewReader("imgUrl="+contents.Param.VerifyImgUrl))
//	}
//	if er != nil {
//		err = base.NewReplayError(base.CODE_PARAM_ERROR, base.CODE_PARAM_ERROR_MESSAGE, nil)
//		return nil, err
//	}
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	resp, er := client.Do(req)
//	defer resp.Body.Close()
//
//	body, er := ioutil.ReadAll(resp.Body)
//	if er != nil {
//		err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
//		return nil, err
//	}
//	fmt.Println(string(body))
//	e := json.Unmarshal(body, &reply)
//	if e != nil {
//		err = base.NewReplayError(base.CODE_UNKNOWN_ERROR, base.CODE_UNKNOWN_ERROR_MESSAGE, nil)
//		return nil, err
//	}
//	contents.Reply.Biz = reply
//	return nil, nil
//
//}
