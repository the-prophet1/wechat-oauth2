package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	langCN  = "zh_CN"
	langTW  = "zh_TW"
	langEng = "en"

	scopeBase     = "snsapi_base"
	scopeUserInfo = "snsapi_userinfo"
)

var (
	oauth2AccessTokenUrl   = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=%s"
	userInfoUrl            = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s"
	weChatAuthorizationUrl = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
)

type Client struct {
	client http.Client
}

type accessTokenByError struct {
	AccessToken `json:",inline"`
	// Code无效错误：错误码
	ErrCode int `json:"errcode"`
	// Code无效错误：错误信息
	ErrMsg string `json:"errmsg"`
}

//AccessToken 从微信第三方授权OAuth2中获取得到的access_token凭证结构
type AccessToken struct {
	//网页授权接口调用凭证
	Token string `json:"access_token"`
	//access_token接口调用凭证超时时间，单位（秒）
	Expires int `json:"expires_in"`
	//用户刷新access_token
	Refresh string `json:"refresh_token"`
	//用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	OpenID string `json:"openid"`
	//用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
}

type userInfoByError struct {
	UserInfo `json:",inline"`
	// Code无效错误：错误码
	ErrCode int `json:"errcode"`
	// Code无效错误：错误信息
	ErrMsg string `json:"errmsg"`
}

type UserInfo struct {
	//用户的唯一标识
	OpenID string `json:"openid"`
	//用户昵称
	NickName string `json:"nickname"`
	//用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Sex int `json:"sex"`
	//用户个人资料填写的省份
	Province string `json:"province"`
	//普通用户个人资料填写的城市
	City string `json:"city"`
	//国家，如中国为CN
	Country string `json:"country"`
	//用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像）
	//用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImgUrl string `json:"headimgurl"`
	//用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	Privilege []string `json:"privilege"`
	//只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段
	UnionID string `json:"unionid"`
}

//GenerateUserInfoAuthorizationUrl 生成用户信息授权地址url
func GenerateUserInfoAuthorizationUrl(appID, redirectUri, state string) string {
	return generateWeChatAuthorizationUrl(appID, redirectUri, scopeUserInfo, state)
}

//GenerateBaseAuthorizationUrl 生成base信息授权地址url
func GenerateBaseAuthorizationUrl(appID, redirectUri, state string) string {
	return fmt.Sprintf(weChatAuthorizationUrl, appID, redirectUri, scopeBase, state)
}

func generateWeChatAuthorizationUrl(appID, redirectUri, scope, state string) string {
	return fmt.Sprintf(weChatAuthorizationUrl, appID, url.QueryEscape(redirectUri), scope, state)
}

//NewClient 创建一个获取网页accessToken的微信客户端
func NewClient() *Client {
	return &Client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (cli *Client) getOauth2AccessTokenUrl(appid, secret, code string) string {
	return fmt.Sprintf(oauth2AccessTokenUrl, appid, secret, code, "authorization_code")
}

func (cli *Client) getOauth2AccessTokenRequest(appid, secret, code string) (*http.Request, error) {
	url := cli.getOauth2AccessTokenUrl(appid, secret, code)
	return http.NewRequest(http.MethodGet, url, nil)
}

func (cli *Client) getUserInfoUrl(accessToken, openid string) string {
	return fmt.Sprintf(userInfoUrl, accessToken, openid, langCN)
}

func (cli *Client) getUserInfoRequest(token *AccessToken) (*http.Request, error) {
	url := cli.getUserInfoUrl(token.Token, token.OpenID)
	return http.NewRequest(http.MethodGet, url, nil)
}

//GetWeChatAccessToken 根据appid,secret,code获取网页授权的accessToken
//appid  公众号的应用id
//secret 与appid相关联的密钥
//code   公众号授权给第三方登陆的code
func (cli *Client) GetWeChatAccessToken(appid, secret, code string) (*AccessToken, error) {
	request, err := cli.getOauth2AccessTokenRequest(appid, secret, code)
	if err != nil {
		return nil, err
	}

	resp, err := cli.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res accessTokenByError
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	} else if res.ErrCode != 0 || len(res.ErrMsg) != 0 {
		return nil, errors.New(res.ErrMsg)
	}

	return &res.AccessToken, nil
}

//GetUserInfoByToken 根据accessToken获取微信的用户信息
//token 获取AccessToken详见GetWeChatAccessToken函数
func (cli *Client) GetUserInfoByToken(token *AccessToken) (*UserInfo, error) {
	req, err := cli.getUserInfoRequest(token)
	if err != nil {
		return nil, err
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res userInfoByError
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	} else if res.ErrCode != 0 || len(res.ErrMsg) != 0 {
		return nil, errors.New(res.ErrMsg)
	}
	return &res.UserInfo, nil
}

//GetUserInfo 根据appid, secret, code获取微信的用户信息
//appid, secret, code 的参数将调用GetWeChatAccessToken获取accessToken
//并将accessToken提交给GetUserInfoByToken用以获取用户信息
func (cli *Client) GetUserInfo(appid, secret, code string) (*UserInfo, error) {
	accessToken, err := cli.GetWeChatAccessToken(appid, secret, code)
	if err != nil {
		return nil, err
	}
	return cli.GetUserInfoByToken(accessToken)
}
