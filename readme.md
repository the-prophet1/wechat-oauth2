# wechat-oauth2

一个微信的网页第三方oauth2授权的客户端

用法:
```go
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

```