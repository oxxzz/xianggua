package origin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"yuewen/store/db"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var client = resty.NewWithClient(&http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	},
})

type oResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func sign() string {
	today := time.Now().Format("20060102")
	h := md5.New()
	h.Write([]byte(today + viper.GetString("secret.ak") + viper.GetString("secret.sk")))
	
	return hex.EncodeToString(h.Sum(nil))
}

type Book struct {
	ID        string `json:"bookId"`
	Name      string `json:"book_name"`
	UpdatedAt string `json:"update_time"`
}

// Books returns all book id list
func Books() ([]*Book, error) {
	url := fmt.Sprintf("http://www.xiangguayuedu.cn/apis/api/BookList.php?sid=%s&sign=%s", viper.GetString("secret.ak"), sign())

	var r struct {
		oResp 
		Data []*Book `json:"data"`
	}
	if _, err := client.R().SetResult(&r).Get(url); err != nil {
		return nil, errors.Wrapf(err, "[ORIGIN.Books] send request failed")
	}

	if r.Code != 200 {
		return nil, fmt.Errorf("[ORIGIN.Books] invalid status, code: %d, msg: %s", r.Code, r.Msg)
	}

	return r.Data, nil
}

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

// Info returns book info by id
func Info(id string) (*BInfo, error) {
	url := "http://www.xiangguayuedu.cn/apis/api/BookInfo.php?sid=%s&sign=%s&bookid=%s"
	var r struct {
		oResp
		Data *BInfo `json:"data"`
	}

	if _, err := client.R().SetResult(&r).Get(fmt.Sprintf(url, viper.GetString("secret.ak"), sign(), id)); err != nil {
		return nil, errors.Wrapf(err, "[ORIGIN.Info] send request failed")
	}

	if r.Code != 200 {
		return nil, fmt.Errorf("[ORIGIN.Info] invalid status. code: %d, msg: %s", r.Code, r.Msg)
	}

	return r.Data, nil
}

type Chapter struct {
	ID        string `json:"chapter_id"`
	// ID on YUEWEN side, used to update chapter
	YWID      string 
	Name      string `json:"chapter_name"`
	UpdatedAt string `json:"update_time"`
	// Content is not returned by API, filled in later
	Content   string
	Vip       int    `json:"is_vip"`
}

// Chapters returns all chapters of a book
func Chapters(id string) ([]*Chapter, error) {
	url := fmt.Sprintf(
		"http://www.xiangguayuedu.cn/apis/api/BookChapters.php?sid=%s&sign=%s&bookid=%s", 
		viper.GetString("secret.ak"), sign(), id)

	var r struct {
		oResp
		Data []struct {
			Volume   string     `json:"volume_name"`
			Chapters []*Chapter `json:"chapterlist"`
		} `json:"data"`
	}

	if _, err := client.R().SetResult(&r).Get(url); err != nil {
		return nil, errors.Wrapf(err, "[ORIGIN.Chapters] send request failed")
	}
	
	if r.Code != 200 {
		return nil, fmt.Errorf("[ORIGIN.Chapters] invalid status. code: %d, msg: %s", r.Code, r.Msg)
	}

	
	if len(r.Data) <= 0 {
		return nil, errors.New("[ORIGIN.Chapters] no chapters")
	}

	items := r.Data[0].Chapters
	for _, item := range items {
		c, err := content(id, item.ID);
		if err != nil {
			return nil, errors.Wrapf(err, "[ORIGIN.Chapters.Content] get chapter %s content failed", item.ID)
		}

		item.Content = c
	}

	return items, nil
}

type CInfo struct {
	ID      string `json:"id"`
	BID     string `json:"bid"`
	Content string `json:"content"`
}

// ChapterContent returns chapter content
func content(bid string, cid string) (string, error) {
	var r struct {
		oResp
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	url := fmt.Sprintf(
		"http://www.xiangguayuedu.cn/apis/api/BookChapterInfo.php?sid=%s&sign=%s&bookid=%s&chapterid=%s", 
		viper.GetString("secret.ak"), sign(), bid, cid)
	if _, err := client.R().SetResult(&r).Get(url); err != nil {
		return "", errors.Wrapf(err, "[ORIGIN.Content] send request failed")
	}
	
	if r.Code != 200 {
		return "", fmt.Errorf("[ORIGIN.Content] code: %d, msg: %s", r.Code, r.Msg)
	}

	return r.Data.Content, nil
}

// has book record
func HasBookRecord(id string) (bool, error) {
	var has bool
	return has, db.MySQL.Get(&has, 
		"select count(1) from yw_books where book_id=?", cast.ToInt(id))
}

func HasChapterRecord(id string, cid string) (bool, error) {
	var has bool
	return has, db.MySQL.Get(&has, 
		"select count(1) from yw_chapters where book_id=? and chapter_id=?", cast.ToInt64(id), cast.ToInt64(cid))
}

func PutBookRecord(b db.YWBook) error {
	tx := db.MySQL.MustBegin()
	defer tx.Rollback()
	if _, err := tx.NamedExec(`
		insert into yw_books (book_id, name, status, yw_book_id, yw_cp_id, book_updated_at)
		values (:book_id, :name, :yw_book_id, :yw_cp_id, :status, :book_updated_at)
	`, b); err != nil {
		return errors.Wrapf(err, "[ORIGIN.PutBookRecord] insert book record failed")
	}

	return tx.Commit()
}

func PutChapterRecord(c *db.YWChapter) (bool, error) {
	tx := db.MySQL.MustBegin()
	defer tx.Rollback()
	if _, err := tx.NamedExec(`
		insert into yw_chapters (book_id, chapter_id, name, yw_cp_id, yw_book_id, yw_chapter_id, status, chapter_updated_at)
		values (:book_id, :chapter_id, :name, :yw_cp_id, :yw_book_id, :yw_chapter_id, :status, :chapter_updated_at)
	`, c); err != nil {
		return false, errors.Wrapf(err, "[ORIGIN.PutChapterRecord] insert chapter record failed")
	}

	return true,tx.Commit()
}

func UpdateBookStatus(b string, status int) error {
	tx := db.MySQL.MustBegin()
	defer tx.Rollback()
	if _, err := tx.Exec("update yw_books set status=? where book_id=?", status, b); err != nil {
		return errors.Wrapf(err, "[ORIGIN.UpdateBookStatus] update book status failed")
	}

	return tx.Commit()
}

func UpdateChapterStatus(b string, c string, status int) error {
	tx := db.MySQL.MustBegin()
	defer tx.Rollback()
	if _, err := tx.Exec("update yw_chapters set status=? where book_id=? and chapter_id=?", status, b, c); err != nil {
		return errors.Wrapf(err, "[ORIGIN.UpdateChapterStatus] update chapterstatus failed")
	}

	return tx.Commit()
}