package models

type AuthenticationRequest struct {
	RelyingPartyName       string `json:"relyingPartyName"`
	RelyingPartyUUID       string `json:"relyingPartyUUID"`
	NationalIdentityNumber string `json:"nationalIdentityNumber"`
	PhoneNumber            string `json:"phoneNumber"`
	Hash                   string `json:"hash"`
	HashType               string `json:"hashType"`
	Language               string `json:"language"`
	DisplayText            string `json:"displayText"`
	DisplayTextFormat      string `json:"displayTextFormat"`
}

type AuthenticationResponse struct {
	State     string    `json:"state"`
	Result    string    `json:"result"`
	Signature Signature `json:"signature"`
	Cert      string    `json:"cert"`
}

type Signature struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}
