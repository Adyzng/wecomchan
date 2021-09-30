package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/adyzng/wecomchan/go-scf/consts"
	"github.com/adyzng/wecomchan/go-scf/dal"
	"github.com/adyzng/wecomchan/go-scf/model"
	"github.com/adyzng/wecomchan/go-scf/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/tencentyun/scf-go-lib/events"
)

var (
	httpCli = http.Client{Timeout: 10 * time.Second}
)

func WeComChanService(ctx context.Context, event events.APIGatewayRequest) map[string]interface{} {
	req := parseRequest(event)
	if req.MsgType == "" || req.Content == "" {
		return utils.MakeResp(-1, "param error")
	}
	if req.SendKey != consts.SEND_KEY {
		return utils.MakeResp(-1, "sendkey error")
	}
	if req.ToUser == "" {
		req.ToUser = consts.WECOM_TOUID
	}
	if err := postWechatMsg(dal.AccessToken, req); err != nil {
		return utils.MakeResp(0, err.Error())
	}
	return utils.MakeResp(0, "success")
}

func postWechatMsg(accessToken string, msg *model.WeComRequest) error {
	content := &model.WechatMsg{
		AgentId:                consts.WECOM_AID,
		ToUser:                 msg.ToUser,
		MsgType:                msg.MsgType,
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 600,
	}

	switch msg.MsgType {
	case "text":
		content.Text = &model.TextMsg{Content: msg.Content}
	case "textcard":
		content.TextCard = &model.TextCardMsg{
			Title:  msg.Title,
			Desc:   msg.Content,
			URL:    msg.JumpURL,
			BtnTxt: "去预约",
		}
	}

	b, _ := jsoniter.Marshal(content)
	req, _ := http.NewRequest("POST", fmt.Sprintf(consts.WeComMsgSendURL, accessToken), bytes.NewBuffer(b))
	req.Header.Set("Content-type", "application/json")

	resp, err := httpCli.Do(req)
	if err != nil {
		fmt.Println("[postWechatMsg] failed, err=", err)
		return nil
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("postWechatMsg statusCode is not 200")
		return errors.New("statusCode is not 200")
	}

	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	postResp := &model.PostResp{}
	if err := jsoniter.Unmarshal(respBodyBytes, postResp); err != nil {
		fmt.Println("postWechatMsg json Unmarshal failed, err=", err)
		return err
	}

	if postResp.Errcode != 0 {
		fmt.Println("postWechatMsg postResp.Errcode != 0, err=", postResp.Errmsg)
		return errors.New(postResp.Errmsg)
	}
	return nil
}

func getQuery(key string, event events.APIGatewayRequest) string {
	switch event.Method {
	case "GET":
		value := event.QueryString[key]
		if len(value) > 0 && value[0] != "" {
			return value[0]
		}
		return ""
	case "POST":
		return jsoniter.Get([]byte(event.Body), key).ToString()
	default:
		return ""
	}
}

func parseRequest(event events.APIGatewayRequest) *model.WeComRequest {
	wxm := &model.WeComRequest{}
	switch event.Method {
	case "GET":
		if tmp := event.QueryString["send_key"]; len(tmp) > 0 {
			wxm.MsgType = tmp[0]
		}
		if tmp := event.QueryString["msg_type"]; len(tmp) > 0 {
			wxm.MsgType = tmp[0]
		}
		if tmp := event.QueryString["title"]; len(tmp) > 0 {
			wxm.Title = tmp[0]
		}
		if tmp := event.QueryString["content"]; len(tmp) > 0 {
			wxm.Content = tmp[0]
		}
		if tmp := event.QueryString["jump_url"]; len(tmp) > 0 {
			wxm.JumpURL = tmp[0]
		}
		if tmp := event.QueryString["to_user"]; len(tmp) > 0 {
			wxm.ToUser = tmp[0]
		}
	case "POST":
		_ = jsoniter.Unmarshal([]byte(event.Body), wxm)
	}
	return wxm
}
