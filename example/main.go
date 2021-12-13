package main

import (
	"fmt"
	wechat "github.com/the-prophet1/wechat-oauth2"
)

var (
	appid  = "xxxxxx"
	secret = "xxxxxx"
	code   = "xxxxxx"
)

func main() {
	cli := wechat.NewClient()

	info, err := cli.GetUserInfo(appid, secret, code)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}
