package controllers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"wuhuaguo.com/whgv01/app/models"
)

type WxApp struct {
	Application
}

func (c WxApp) Wx(signature string,
	timestamp string, nonce string,
	echostr string) revel.Result {
	revel.INFO.Println(signature, timestamp, nonce, echostr)
	// token := "1987526ab"
	return c.RenderText(echostr)
}

func (c WxApp) WxP() revel.Result {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
		return c.RenderText("error")
	}

	str := string(body)
	revel.INFO.Println(str)
	msg := models.ParseMsg(str)

	tu := &models.CDDATA{V: msg.FromUserName}
	fu := &models.CDDATA{V: msg.ToUserName}
	mt := &models.CDDATA{V: "text"}
	ct := &models.CDDATA{V: "resp: " + msg.Content}

	rm := &models.RspMsg{
		ToUserName:   tu,
		FromUserName: fu,
		CreateTime:   msg.CreateTime,
		MsgType:      mt,
		Content:      ct,
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent("  ", "    ")
	if err := enc.Encode(rm); err != nil {
		revel.ERROR.Println("error: %v\n", err)
	}
	rsmsg := b.String()
	revel.INFO.Println(rsmsg)
	return c.RenderText(rsmsg)
}

func (c WxApp) Login(code string, state string) revel.Result {
	revel.INFO.Println("code %s state %s", code, state)
	if code == "" {
		return c.RenderText("login fail %q", state)
	} else {
		//获取到了用户id
		resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=wx82c3aa347250de4b&secret=60942365de45e1297cc50d7ead3371d4&code=%s&grant_type=authorization_code", code))
		if err != nil {
			return c.RenderText(err.Error())
		} else {
			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				return c.RenderText(err1.Error())
			}
			str := string(body)
			//解析json
			wa := models.ParseWxWebAccessToken(str)
			revel.INFO.Println(str, (*wa).Openid)
			if wa != nil {
				return c.RenderText("%q has login", (*wa).Openid)
			} else {
				return c.RenderText("获取access_token 错误")
			}
		}
	}

}

func (c WxApp) Index() revel.Result {
	user := c.wxId()
	if user != nil {
		revel.INFO.Println("WxApp/Index: " + user.Name + " " + c.Request.UserAgent())
	} else {
		return c.Redirect("https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx82c3aa347250de4b&redirect_uri=http%3a%2f%2f123.207.143.158%2fwxlogin&response_type=code&scope=snsapi_base&state=123#wechat_redirect")
	}
	//获取boxlist 商品
	irs, err := c.boxListItems()
	if err != nil {
		revel.ERROR.Println(err)
	}
	return c.Render(user, irs)
}
