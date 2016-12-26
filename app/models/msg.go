package models

import (
	"encoding/xml"
	"github.com/revel/revel"
)

type PushMsg struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        string `xml:"MsgId"`
}

func ParseMsg(str string) *PushMsg {
	v := &PushMsg{}
	err := xml.Unmarshal([]byte(str), v)
	if err != nil {
		revel.ERROR.Println(err, "parse push message error")
		return nil
	}
	return v
}

type CDDATA struct {
	V string `xml:",cdata"`
}

type RspMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   *CDDATA
	FromUserName *CDDATA
	CreateTime   string
	MsgType      *CDDATA
	Content      *CDDATA
}
