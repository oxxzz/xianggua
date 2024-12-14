package db

import "time"

type YWBook struct {
	ID            int64      `json:"id" db:"id"`
	BookID        string     `json:"book_id" db:"book_id"`
	Name          string     `json:"name" db:"name"`
	YWCPID        int64      `json:"yw_cp_id" db:"yw_cp_id"`
	YWBookID      string     `json:"yw_book_id" db:"yw_book_id"`
	Status        int        `json:"status" db:"status"`
	BookUpdatedAt int64      `json:"book_updated_at" db:"book_updated_at"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

type YWChapter struct {
	ID               int64      `json:"id" db:"id"`
	BookID           string     `json:"book_id" db:"book_id"`
	Name             string     `json:"name" db:"name"`
	ChapterID        string     `json:"chapter_id" db:"chapter_id"`
	YWCPID           int64      `json:"yw_cp_id" db:"yw_cp_id"`
	YWBookID         string     `json:"yw_book_id" db:"yw_book_id"`
	YWChapterID      string     `json:"yw_chapter_id" db:"yw_chapter_id"`
	Status           int        `json:"status" db:"status"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`
	ChapterUpdatedAt int64      `json:"chapter_updated_at" db:"chapter_updated_at"`
}
