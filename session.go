package mobileid

type Session struct {
	Id   string `json:"sessionID"`
	Code string `json:"code"`
}
