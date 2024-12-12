package origin

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	},
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func sign() string {
	today := time.Now().Format("20060102")
	h := md5.New()
	h.Write([]byte(today + viper.GetString("secret.ak") + viper.GetString("secret.sk")))
	s := hex.EncodeToString(h.Sum(nil))

	return s
}

//	{
//	  "bookId": "282",
//	  "book_name": "试传文件22",
//	  "update_time": "1510992075"
//	}
type Book struct {
	ID        string `json:"bookId"`
	Name      string `json:"book_name"`
	UpdatedAt string `json:"update_time"`
}

func Books() ([]*Book, error) {
	url := fmt.Sprintf("http://www.xiangguayuedu.cn/apis/api/BookList.php?sid=%s&sign=%s", viper.GetString("secret.ak"), sign())
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var resp struct {
		resp
		Data []*Book `json:"data"`
	}

	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, msg: %s", resp.Code, resp.Msg)
	}

	return resp.Data, nil
}

// "bookId": "61",
// "book_name": "这是一本测试小说",
// "author": "测试作者",
// "brief": "测试简介”",
// "words": "710092",
// "keywords": "",
// "cover": "http://www.xiangguayuedu.cn/images/nocover.jpg",
// "group_id": "1008",
// "cate_id": "1008002",
// "is_vip": 1,
// "update_time": "1542770256",
// "status": 1
type BInfo struct {
	ID       string `json:"bookId"`
	Name     string `json:"book_name"`
	Author   string `json:"author"`
	Brief    string `json:"brief"`
	Words    string `json:"words"`
	Cover    string `json:"cover"`
	GroupID  string `json:"group_id"`
	CateID   string `json:"cate_id"`
	Vip      int    `json:"is_vip"`
	UpdateAt string `json:"update_time"`
	Status   int    `json:"status"`
}

func Info(id string) (*BInfo, error) {
	url := "http://www.xiangguayuedu.cn/apis/api/BookInfo.php?sid=%s&sign=%s&bookid=%s"
	r, err := client.Get(fmt.Sprintf(url, viper.GetString("secret.ak"), sign(), id))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var resp struct {
		resp
		Data *BInfo `json:"data"`
	}

	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, msg: %s", resp.Code, resp.Msg)
	}

	return resp.Data, nil
}

// {
//   "chapter_id": "4012",
//   "chapter_name": "第一章 故人1",
//   "update_time": "1542770186",
//   "is_vip": 0
// }

type Chapter struct {
	ID        string `json:"chapter_id"`
	Name      string `json:"chapter_name"`
	UpdatedAt string `json:"update_time"`
	Vip       int    `json:"is_vip"`
}

func Chapters(id string) ([]*Chapter, error) {
	url := "http://www.xiangguayuedu.cn/apis/api/BookChapters.php?sid=%s&sign=%s&bookid=%s"
	r, err := client.Get(fmt.Sprintf(url, viper.GetString("secret.ak"), sign(), id))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var resp struct {
		resp
		Data []struct {
			Volume   string     `json:"volume_name"`
			Chapters []*Chapter `json:"chapterlist"`
		} `json:"data"`
	}

	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, msg: %s", resp.Code, resp.Msg)
	}

	if len(resp.Data) > 0 {
		return resp.Data[0].Chapters, nil
	}

	return nil, errors.New("no chapters")
}

// {
//   "code": 200,
//   "data": {
//     "content": "　　“以前他还会脸红呢。”罗娜感叹道，“青春一去不复返啊。”\n　　吴泽冷笑了一声。\n　　两周后，段宇成出发参加田径大奖赛第一站。"
//   },
//   "msg": "Success"
// }

type CInfo struct {
	ID      string `json:"id"`
	BID     string `json:"bid"`
	Content string `json:"content"`
}

func ChapterInfo(bid string, cid string) (*CInfo, error) {
	url := "http://www.xiangguayuedu.cn/apis/api/BookChapterInfo.php?sid=%s&sign=%s&bookid=%s&chapterid=%s"
	r, err := client.Get(fmt.Sprintf(url, viper.GetString("secret.ak"), sign(), bid, cid))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var resp struct {
		resp
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("code: %d, msg: %s", resp.Code, resp.Msg)
	}

	return &CInfo{
		Content: resp.Data.Content,
		ID:      cid,
		BID:     bid,
	}, nil
}
