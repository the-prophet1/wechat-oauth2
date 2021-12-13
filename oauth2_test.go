package wechat

import (
	"encoding/json"
	"testing"
)

var (
	appid  = "wxd31cc88f77409a67"
	secret = "0eea140ef5acb8508b2af4b84563ee2e"
	code   = "011YKpll2l1fg84TWQkl2Ujba61YKplQ"
)

func TestClient_GetUserInfo(t *testing.T) {
	client := NewClient()

	token, err := client.GetWeChatAccessToken(appid, secret, code)
	if err != nil {
		t.Fatal("get accessToken error", err)
	}

	tokenStruct, _ := json.Marshal(token)
	t.Log("token", string(tokenStruct))
}

func TestClient_GetWeChatAccessToken(t *testing.T) {
	client := NewClient()
	info, err := client.GetUserInfo(appid, secret, code)
	if err != nil {
		t.Fatal("get userInfo error", err)
	}

	infoStruct, _ := json.Marshal(info)
	t.Log("info", string(infoStruct))
}

func TestClient_GetUserInfoByToken(t *testing.T) {
	client := NewClient()
	token, err := client.GetWeChatAccessToken(appid, secret, code)
	if err != nil {
		t.Fatal("get accessToken error", err)
	}
	tokenStruct, _ := json.Marshal(token)
	t.Log("token", string(tokenStruct))

	info, err := client.GetUserInfoByToken(token)
	if err != nil {
		t.Fatal("get userInfo error", err)
	}
	infoStruct, _ := json.Marshal(info)
	t.Log("info", string(infoStruct))
}

func TestGenerateBaseAuthorizationUrl(t *testing.T) {
	url := GenerateBaseAuthorizationUrl(
		"sn123fdsbnj934uq13",
		"https://www.baidu.com/wechat/authorize",
		"1",
	)

	t.Log(url)
}

func TestGenerateUserInfoAuthorizationUrl(t *testing.T) {
	url := GenerateUserInfoAuthorizationUrl(
		"sn123fdsbnj934uq13",
		"https://www.ykkcfw.com/wechat/authorize",
		"1",
	)

	t.Log(url)
}
