package discord

type Message struct {
	Content             string `json:"content"`
	Flags               int64  `json:"flags"`
	Mobile_network_type string `json:"mobile_network_type"`
	Nonce               string `json:"nonce"`
	Tts                 bool   `json:"tts"`
}

type Response struct {
	Msg  string `json:"message"`
	Code int64  `json:"code"`
}
