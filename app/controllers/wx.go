package controllers

import (
	"github.com/revel/revel"
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
