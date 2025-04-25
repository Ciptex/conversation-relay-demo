package llmhandlers

import (
	"bytes"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"html/template"
)

type prompt struct {
	callSid     string
	lastMessage string
	accountSid  string
	configSid   string
	repo        repo.IRepo
	span        trace.ISpan
}

func newPrompt(callSid, accountSid, configSid string, repo repo.IRepo, span trace.ISpan) *prompt {
	return &prompt{
		callSid:    callSid,
		accountSid: accountSid,
		configSid:  configSid,
		repo:       repo,
		span:       span,
	}
}

func parseTemplate(templateStr, templateName string, templateVars map[string]interface{}) (string, int, error) {
	t, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", 0, err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, templateVars)
	if err != nil {
		return "", 0, err
	}
	tmpStr := tpl.String()
	return tmpStr, 0, nil
}

func (p *prompt) getGenericPrompt(template string) (string, error) {
	sc := (map[string]interface{}{
		"message": p.lastMessage,
	})
	parsedTemplate, _, err := parseTemplate(template, "generic_prompt", sc)
	return parsedTemplate, err
}
