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
	
	if _, err := client.R().SetBody(param.Encode()).SetResult(&r).Post(fmt.Sprintf(Login, u, secret)); err != nil {
		return "", err
	}

	if r.Code != 0 {
		return "", fmt.Errorf("[YUEWEN.Login] invalid status. code: %d, msg: %s", r.Code, r.Message)
	}

	return r.Result.Key, nil
}

// 2 Push book
func PushBookInfo(b *origin.BInfo) (string, error) {
	key, err := SignIn(false)
	if err != nil {
		return "", errors.Wrapf(err, "[YUEWEN.PushBookInfo] login failed")
	}

	free := 1
	if b.Vip == 1 {
		free = 0
	}

	p := map[string]string{
		"key": key,
		"b.cpid": viper.GetString("yuewen.cpid"),
		"b.cpBid": b.ID,
		"b.title ": b.Name,
		"b.author ": b.Author,
		"b.finish": cast.ToString(b.Status),
		"b.intro ": b.Brief,
		"b.free ": cast.ToString(free),
		"b.language": "0",
		"b.form": "0",
		// "b.sex ": "性别属性", // ?
		// "b.maxFreeChapter": "最大免费章节数",
		// "b.wordPrice": "千字价格",
		// "b.category": "分类",
	}

	var r struct {
		resp
		Result struct {
			ID string `json:"bookid"`
		} `json:"result"`
	}

	cover := strings.Join([]string{os.TempDir(), b.ID + ".jpg"}, "")
	if _, err := client.R().EnableTrace().SetOutput(cover).Get(b.Cover); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(
			"[YUEWEN.PushBookInfo] book: %s, cover: %s, download failed", b.ID, b.Cover))
	}

	if _, err := os.Stat(cover); errors.Is(err, os.ErrNotExist) {
		return "", errors.Wrap(err, fmt.Sprintf(
			"[YUEWEN.PushBookInfo] book: %s, cover: %s, download failed", b.ID, b.Cover))
	} 

	defer os.RemoveAll(cover)
	if _, err = client.R().SetFormData(p).SetFile("b.cover", cover).SetResult(&r).Post(PushBook) ;err != nil {
		return "", errors.Wrapf(err, "[YUEWEN.PushBookInfo] push book %s failed", b.ID)
	}

	if r.Code != 0 {
		return "", fmt.Errorf("[YUEWEN.PushBookInfo] invalid status, code: %d, msg: %s", r.Code, r.Message)
	}

	return r.Result.ID, nil
}

// 2.1 Push chapter
func PushChapter(bid string, c *origin.Chapter) error {
	key, err := SignIn(false)
	if err != nil {
		return errors.Wrapf(err, "[YUEWEN.PushChapter] login failed")
	}

	p := map[string]interface{}{
		"key": key,
		"c.cpid": viper.GetInt64("yuewen.cpid"),
		"c.bookid": cast.ToInt64(bid),
		"c.title ": c.Name,
		"c.content ": c.Content,
		"c.cpcid": c.ID,
	}

	var r struct {
		resp
		Result struct {
			ID string `json:"bookid"`
		} `json:"result"`
	}

	if _, err = client.R().SetBody(p).SetResult(&r).Post(PushChaoter); err != nil {
		return errors.Wrapf(err, "[YUEWEN.PushChapter] push chapter %s.%s failed", bid, c.ID)
	}

	if r.Code != 0 {
		return fmt.Errorf("[YUEWEN.PushChapter] code: %d, msg: %s", r.Code, r.Message)
	}

	return nil
}
// 2.2 Fetch chapter status
// 2.3 Update chapter
func PatchChapter(bid string, c *origin.Chapter) error {
	key, err := SignIn(false)
	if err != nil {
		return errors.Wrapf(err, "[YUEWEN.PatchChapter] login failed")
	}

	p := map[string]interface{}{
		"key": key,
		"c.cpid": viper.GetInt64("yuewen.cpid"),
		"c.bookid": cast.ToInt64(bid),
		"c.title ": c.Name,
		"c.content ": c.Content,
		"c.chapterid": c.YWID,
		"c.cpcid": c.ID,
	}

	var r struct {
		resp
		Result struct {
			ID string `json:"bookid"`
		} `json:"result"`
	}

	if _, err = client.R().SetBody(p).SetResult(&r).Post(PushChaoter); err != nil {
		return errors.Wrapf(err, "[YUEWEN.PatchChapter] push chapter %s.%s failed", bid, c.ID)
	}

	if r.Code != 0 {
		return fmt.Errorf("[YUEWEN.PatchChapter] code: %d, msg: %s", r.Code, r.Message)
	}

	return nil
}
