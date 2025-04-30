package twilio

import (
	"bytes"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/twilio/twilio-go/twiml"
)

type Twiml struct {
	span trace.ISpan
}

func NewTwiml(span trace.ISpan) *Twiml {
	return &Twiml{
		span: span,
	}
}

// var twiml_1 = `<Response><Connect action="https://webhooks.twilio.com/v1/Accounts/{{.aAccsid}}/Flows/{{.flowSid}}"><ConversationRelay url="wss://{{.url}}"><Parameter name="accSid" value="{{.accSid}}"/><Parameter name="configSid" value="{{.configSid}}"/></ConversationRelay></Connect></Response>`
var twiml_1 = `<Response><Connect action="https://{{.actionUrl}}"><ConversationRelay url="wss://{{.url}}"><Parameter name="accSid" value="{{.accSid}}"/><Parameter name="configSid" value="{{.configSid}}"/></ConversationRelay></Connect></Response>`

func (t *Twiml) CreateConversationRelayTwiml(accSid, congfigSid, host string) (string, error) {
	sc := (map[string]interface{}{
		"accSid":    accSid,
		"configSid": congfigSid,
		"url":       host + "/v1.0/" + congfigSid + "/ws",
		"actionUrl": fmt.Sprintf("%s/v1.0/%s/%s/queue", host, accSid, congfigSid),
	})
	xmlHeader := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	twimlStr, err := parseTemplate(twiml_1, "twiml", sc)
	return xmlHeader + twimlStr, err
}

func parseTemplate(templateStr, templateName string, templateVars map[string]interface{}) (string, error) {
	t, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, templateVars)
	if err != nil {
		return "", err
	}
	tmpStr := tpl.String()
	return tmpStr, nil
}

func (t *Twiml) EnqueueCall(config types.AccountConfig, host string) string {
	q := &twiml.VoiceEnqueue{
		Name:        "support",
		WorkflowSid: config.TwilioWorkFlowSid,
	}
	attrs, _ := json.Marshal(map[string]string{
		"account":       "123",
		"paymentMethod": "credit card",
		"amount":        "100",
	})
	q.InnerElements = append(q.InnerElements, &twiml.VoiceTask{
		Body: string(attrs),
	})
	verbList := []twiml.Element{q}
	twimlResult, _ := twiml.Voice(verbList)
	return twimlResult
}

// func (t *Twiml) UpdateCall(accConfig types.AccountConfig, callSid string) error {
// 	client := createTwilioClient(accConfig)
// 	redirect := &twiml.VoiceRedirect{
// 		Url:    fmt.Sprintf("https://webhooks.twilio.com/v1/Accounts/%s/Flows/%s", accConfig.AccountSid, accConfig.TwilioFlowSid),
// 		Method: "POST",
// 	}
// 	_ = redirect
// 	// say := &twiml.VoiceSay{
// 	// 	Message: "hello how are you",
// 	// }
// 	verbList := []twiml.Element{redirect}
// 	twimlResult, err := twiml.Voice(verbList)
// 	if err == nil {
// 		t.span.Debug("Twilio::UpdateCall::", "twiml", twimlResult)
// 	} else {
// 		t.span.Error("Twilio::UpdateCall::", err)
// 		return err
// 	}
// 	params := &twilioApi.UpdateCallParams{
// 		PathAccountSid: &accConfig.AccountSid,
// 		Twiml:          &twimlResult,
// 	}
// 	_, err = client.Api.UpdateCall(callSid, params)
// 	if err == nil {
// 		t.span.Debug("twilio::Say::updated call", "sid", callSid)
// 	} else {
// 		t.span.Error("twilio::Say::error updating call", err)
// 	}
// 	return err
// }
