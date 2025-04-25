package twilio

import (
	"bytes"
	"html/template"
	"os"
)

type Twiml struct {
}

func NewTwiml() *Twiml {
	return &Twiml{}
}

var twiml_1 = `<Response><Connect><ConversationRelay url="wss://{{.url}}"><Parameter name="accSid" value="{{.accSid}}"/><Parameter name="configSid" value="{{.configSid}}"/></ConversationRelay></Connect></Response>`

func (t *Twiml) CreateConversationRelayTwiml(accSid, congfigSid string) (string, error) {
	domain := os.Getenv("PUBLIC_ENDPOINT")
	sc := (map[string]interface{}{
		"accSid":    accSid,
		"configSid": congfigSid,
		"url":       domain + "/v1.0/" + congfigSid + "/ws",
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
