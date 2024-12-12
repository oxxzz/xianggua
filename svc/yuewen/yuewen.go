package yuewen

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	Login         = "https://open.book.qq.com/login?username=%s&password=%s"
	PushBook      = "https://open.book.qq.com/push/addCofreeBook"
	PushChaoter   = "https://open.book.qq.com/push/addChapter"
	FetchChapter  = "https://open.book.qq.com/push/getUpdateInfo"
	UpdateChapter = "https://open.book.qq.com/push/updateChapter"
)

// 1 Login
// 密码 = md5(md5(cpid)+md5(平台登录密码)+md5(username))
func SignIn(debug bool) (string, error) {
	u := viper.GetString("yuewen.user")
	p := viper.GetString("yuewen.secret")
	c := viper.GetString("yuewen.cpid")
	pSegments := []string{c, p, u}
	fmt.Println(pSegments)
	_ = pSegments

	return "", nil
}

// 2 Push book
// 2.1 Push chapter
// 2.2 Fetch chapter status
// 2.3 Update chapter
