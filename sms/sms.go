package sms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"wangqingang/sms_lib"
	"wangqingang/sms_lib/pb"

	"wangqingang/cunxun/common"
	e "wangqingang/cunxun/error"
)

const (
	aliSignName         = "大脸"
	aliTemplateCode     = "SMS_89765057"
	aliTemplateParamFmt = "{\"code\": \"%s\"}"
)

func SendCheckcode(config *common.SmsConfig, phone, purpose, checkcode string) (string, error) {
	request := pb.SendSmsRequest{
		PhoneNumbers:  phone,
		SignName:      aliSignName,
		TemplateCode:  aliTemplateCode,
		TemplateParam: fmt.Sprintf(aliTemplateParamFmt, checkcode),
	}

	if purpose != common.SignupPurpose {
		return "", e.S(e.MParamsErr, e.ParamsInvalidPurpose)
	}
	url := sms_lib.CreateSmsSendUrlWithAccess(config.AliAccessId, config.AliAccessSecret, &request)
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return "", e.SP(e.MSmsErr, e.SmsSendErr, err)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", e.SP(e.MSmsErr, e.SmsReadResponse, err)
	}
	smsResponse := &pb.SendSmsRespose{}
	if err := json.Unmarshal(bytes, smsResponse); err != nil {
		return "", e.SP(e.MSmsErr, e.SmsDecodeResponse, err)
	}

	return string(bytes), nil
}
