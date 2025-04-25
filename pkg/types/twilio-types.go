package types

// Inbound payload from Twilio's conversation relay
type TwilioCRInboundPayload struct {
	Type                      string            `json:"type"`
	SessionId                 string            `json:"sessionId"`
	CallSid                   string            `json:"callSid"`
	ParentCallSid             string            `json:"parentCallSid"`
	From                      string            `json:"from"`
	To                        string            `json:"to"`
	ForwardedFrom             string            `json:"forwardedFrom"`
	CallerName                string            `json:"callerName"`
	Direction                 string            `json:"direction"`
	CallType                  string            `json:"callType"`
	CallStatus                string            `json:"callStatus"`
	AccountSid                string            `json:"accountSid"`
	ApplicationSid            string            `json:"applicationSid"`
	VoicePrompt               string            `json:"voicePrompt"`
	Lang                      string            `json:"lang"`
	Last                      bool              `json:"last"`
	CustomParameters          map[string]string `json:"customParameters"`
}
