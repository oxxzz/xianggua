package yuewen

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"yuewen/store/svc/origin"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	Login         = "https://open.book.qq.com/push/login?username=%s&password=%s"
	PushBook      = "https://open.book.qq.com/push/addCofreeBook"
	PushChaoter   = "https://open.book.qq.com/push/addChapter"
	FetchChapter  = "https://open.book.qq.com/push/getUpdateInfo"
	UpdateChapter = "https://open.book.qq.com/push/updateChapter"
)

type resp struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

var lock sync.Mutex
var key = ""
var client = resty.NewWithClient(&http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 200,
	},
})

// 1 Login
// 密码 = md5(md5(cpid)+md5(平台登录密码)+md5(username))
func SignIn(force bool) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	if key != "" && !force {
		return key, nil
	}

	u := viper.GetString("yuewen.user")
	p := viper.GetString("yuewen.secret")
	c := viper.GetString("yuewen.cpid")
	pSegments := []string{c, p, u}
	pString := ""
	h := md5.New()
	for _, v := range pSegments {
		h.Reset()
		h.Write([]byte(v))
		pString += hex.EncodeToString(h.Sum(nil))
	}

	h.Reset()
	h.Write([]byte(pString))
	secret := hex.EncodeToString(h.Sum(nil))

	param := url.Values{}
	param.Add("username", u)
	param.Add("password", secret)
	
	var r struct {
		resp
		Result struct {
			Key string `json:"key"`
		} `json:"result"`
	}
	
	_, err := client.R().SetBody(param.Encode()).SetResult(&r).Post(fmt.Sprintf(Login, u, secret))
	if err != nil {
		return "", err
	}

	if r.Code != 0 {
		return "", fmt.Errorf("code: %d, msg: %s", r.Code, r.Message)
	}

	key = r.Result.Key
	logrus.WithField("KEY", key).Debugf("[YUEWEN] login success")
	return key, nil
}

// 2 Push book
func PushBookInfo(b *origin.BInfo) (string, error) {
	key, err := SignIn(false)
	if err != nil {
		return "", err
	}

	p := map[string]interface{}{
		"key": key,
		"b.cpid": viper.GetInt64("yuewen.cpid"),
		"b.cpBid": cast.ToInt64(b.ID),
		"b.title ": b.Name,
		"b.author ": b.Author,
		"b.finish": b.Status,
		"b.intro ": b.Brief,
		"b.sex ": "", // ?
		"b.free ": cast.ToInt(b.Vip == 0),
		"b.category": 0,
	}

	var r struct {
		resp
		Result struct {
			ID string `json:"bookid"`
		} `json:"result"`
	}

	cover := strings.Join(strings.Split(os.TempDir(), b.ID + ".jpg"), "")
	_, err = client.R().SetOutput(cover).Get(b.Cover)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(
			"book: %s, cover: %s, download failed", b.ID, b.Cover))
	}

	if _, err := os.Stat(cover); errors.Is(err, os.ErrNotExist) {
		return "", errors.Wrap(err, fmt.Sprintf(
			"book: %s, cover: %s, download failed", b.ID, b.Cover))
	} 

	_, err = client.R().SetBody(p).
	SetFile("b.cover", cover).
	SetResult(&r).Post(PushBook)
	if err != nil {
		return "", errors.Wrapf(err, "push book %s failed: %s", b.ID, err.Error())
	}

	defer os.RemoveAll(cover)
	fmt.Println(r)
	return "", nil
}

// 2.1 Push chapter
// 2.2 Fetch chapter status
// 2.3 Update chapter
