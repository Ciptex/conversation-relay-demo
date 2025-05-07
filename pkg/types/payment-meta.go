package types

type PaymentMeta struct {
	Epid      string `json:"epid"`
	Pid       string `json:"pid"`
	Token     string `json:"token"`
	LuhnToken string `json:"luhnToken"`
}
