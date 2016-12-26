package controllers

import (
	"github.com/revel/revel"
	"io/ioutil"
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
	return c.RenderText("<xml><ToUserName><![CDATA[oKvGywSzMr5dSYtDTTHOQxzSoCT8]]></ToUserName>
							<FromUserName><![CDATA[gh_bacadd9c67fd]]></FromUserName>
							<CreateTime>1482747385</CreateTime>
							<MsgType><![CDATA[text]]></MsgType>
							<Content><![CDATA[123]]></Content>
						</xml>")
}
