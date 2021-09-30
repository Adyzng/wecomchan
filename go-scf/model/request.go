package model

type WeComRequest struct {
	SendKey string `json:"send_key"`
	MsgType string `json:"msg_type"` // wx消息类型
	Title   string `json:"title"`    // 标题
	Content string `json:"content"`  // 消息内容
	JumpURL string `json:"jump_url"` // textcard:调转链接
	ToUser  string `json:"to_user"`
}
