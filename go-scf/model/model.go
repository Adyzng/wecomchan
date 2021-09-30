package model

type AssesTokenResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// https://work.weixin.qq.com/api/doc/90002/90151/90854
type WechatMsg struct {
	ToUser                 string       `json:"touser"`
	AgentId                string       `json:"agentid"`
	MsgType                string       `json:"msgtype"`
	EnableDuplicateCheck   int          `json:"enable_duplicate_check"`
	DuplicateCheckInterval int          `json:"duplicate_check_interval"`
	Text                   *TextMsg     `json:"text,omitempty"`
	TextCard               *TextCardMsg `json:"textcard,omitempty"`
}

type TextMsg struct {
	Content string `json:"content"`
}

type TextCardMsg struct {
	Title  string `json:"title"`
	Desc   string `json:"description"`
	URL    string `json:"url"`
	BtnTxt string `json:"btntxt"`
}

type PostResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	Invaliduser string `json:"invaliduser"`
}
