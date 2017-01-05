package models

import (
	"fmt"
	"testing"
)

func TestParseWxWebAccessToken(t *testing.T) {
	str := `{"access_token":"eIFpnL-e-0gKUPST0FFaWQZq3JfTfREBKx8iqxxvpEOmSuO2sefyWZlNRifcNgk1RTA6oZ27Bia7Zsdcv328Cb6234T35eAS0jMEmYlu0h0","expires_in":7200,"refresh_token":"vkIOsPuPWpNQS_DIUfSYUtow-6WkOTcKwKe4-w_SnmL8kTP1ZYMfV4GbNA0TnFG9HaY1fP5q4pisoJQw8gYNXRIQTdKxcQkiU9Ioy0MXrTQ","openid":"oKvGywSzMr5dSYtDTTHOQxzSoCT8","scope":"snsapi_base"}`
	obj := ParseWxWebAccessToken(str)
	if obj == nil {
		t.Errorf("object is nil")
	}
	fmt.Println(obj.String())
	fmt.Println("nihao")
	if obj.Errcode != 0 {
		t.Errorf("error code %d", obj.Errcode)
	}
	fmt.Println(obj)
	fmt.Println(obj.Openid)
}
