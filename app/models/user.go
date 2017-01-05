package models

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"io"
	"regexp"
	"strings"
)

type User struct {
	UserId             int
	Name               string
	Username, Password string
	HashedPassword     []byte
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}

/*
* 网页授权接口凭证
 */
type WebAccessTokenObj struct {
	access_token  string
	expires_in    int
	refresh_token string
	openid        string
	scope         string
	errcode       int
	errmsg        string
}

func (w *WebAccessTokenObj) Access_token() string {
	return w.access_token
}
func (w *WebAccessTokenObj) Refresh_token() string {
	return w.refresh_token
}
func (w *WebAccessTokenObj) Openid() string {
	return w.openid
}
func (w *WebAccessTokenObj) Scope() string {
	return w.scope
}
func (w *WebAccessTokenObj) Errmsg() string {
	return w.errmsg
}
func (w *WebAccessTokenObj) Errcode() int {
	return w.errcode
}

func ParseWxWebAccessToken(jstr string) *WebAccessTokenObj {
	dec := json.NewDecoder(strings.NewReader(jstr))
	obj := &WebAccessTokenObj{}
	if err := dec.Decode(obj); err == io.EOF {
		return nil
	} else if err != nil {
		revel.ERROR.Println(err.Error())
		return nil
	}
	if obj.Errcode() != 0 {
		revel.ERROR.Printf("get weixin web access_token error %d, msg: %s", obj.Errcode(), jstr)
		return nil
	}

	return obj
}
