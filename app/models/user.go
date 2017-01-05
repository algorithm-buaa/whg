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
	Access_token  string
	Expires_in    int
	Refresh_token string
	Openid        string
	Scope         string
	Errcode       int
	Errmsg        string
}

func (w *WebAccessTokenObj) String() string {
	return w.Access_token + " 123\t" +
		w.Refresh_token + "123\t" +
		w.Openid + "\t" +
		w.Scope
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

	if obj.Errcode != 0 {
		revel.ERROR.Printf("get weixin web access_token error %d, msg: %s", obj.Errcode, jstr)
		return nil
	}

	return obj
}
