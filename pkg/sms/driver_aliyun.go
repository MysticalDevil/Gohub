package sms

import (
	"encoding/json"
	aliyunsmsclient "github.com/KenmyZhang/aliyun-communicate"
	"gohub/pkg/logger"
)

// Aliyun Implement the sms.Driver interface
type Aliyun struct{}

// Send Implement the Send method of the sms.Driver interface
func (s *Aliyun) Send(phone string, message Message, config map[string]string) bool {
	smsClient := aliyunsmsclient.New("http://dysmsapi.aliyuncs.com")

	templateParam, err := json.Marshal(message.Data)
	if err != nil {
		logger.ErrorString("SMS[Aliyun]", "Parse binding error", err.Error())
		return false
	}

	logger.DebugJSON("SMS[Aliyun]", "Config information", config)

	result, err := smsClient.Execute(
		config["access_key_id"],
		config["access_key_secret"],
		phone,
		config["sign_name"],
		message.Template,
		string(templateParam),
	)

	logger.DebugJSON("SMS[Aliyun]", "Request content", smsClient.Request)
	logger.DebugJSON("SMS[Aliyun]", "Interface response", result)

	if err != nil {
		logger.ErrorString("SMS[Aliyun]", "Send SMS failed", err.Error())
		return false
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.ErrorString("SMS[Aliyun]", "Parse response JSON error", err.Error())
		return false
	}

	if result.IsSuccessful() {
		logger.DebugString("SMS[Aliyun]", "Send SMSs succeed", "")
		return true
	} else {
		logger.ErrorString("SMS[Aliyun]", "The service provider returns an error", string(resultJSON))
		return false
	}
}
