package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"../common"
	_ "github.com/lib/pq" // this driver for postgres

	"github.com/grokify/html-strip-tags-go"
)

func Top(w http.ResponseWriter, r *http.Request) {

	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	type BreadCrumb struct {
		Level        string
		CategoryID   int
		CategoryName string
	}
	type Note struct {
		NoteID        int
		UpdatedAt     string
		NoteTitle     string
		NoteTxt       string
	}
	type CategoryList struct {
		Level            string
		CategoryID       int
		CategoryName     string
		// Note             []Note
	}
	type View struct {
		CacheV              string
		CSRF                string
		BreadCrumb          []BreadCrumb
		CategoryList        []CategoryList
		CategoryName        string
		CategoryDescription string
		CategoryTxt         template.HTML
		Note                []Note
	}
	var view View
	view.CacheV = common.CacheV
	view.CSRF = ""

	rows, err := db.Query(`SELECT category_id, category_name FROM m_category_name WHERE category_id in 
		(SELECT level_1 FROM m_category_tree GROUP BY level_1)`)
	if err != nil {
		log.Print(err)
	}
	var categoryList []CategoryList
	for rows.Next() {
		y := CategoryList{}
		if err := rows.Scan(&y.CategoryID, &y.CategoryName); err != nil {
			log.Print(err)
		}
		y.Level = "1"
		categoryList = append(categoryList, y)
	}

	rows, err = db.Query(`SELECT note_id, note_title, note_txt, updated_at
			FROM t_note ORDER BY note_id DESC LIMIT 8`)
	if err != nil {
		log.Print(err)
	}
	var notes []Note
	for rows.Next() {
		var ti time.Time
		var txt string
		r := Note{}
		if err := rows.Scan(&r.NoteID, &r.NoteTitle, &txt, &ti); err != nil {
			log.Print(err)
		}
		// fmt.Printf("r.UpdatedAt %#v\n", ti.Format("2006年1月2日"))
		r.UpdatedAt = ti.Format("2006年1月2日")
		r.NoteTxt = strip.StripTags(txt)[0:240]
		notes = append(notes, r)
	}
	view.CategoryName = "炎上案件上等の社外CTO"
	view.CategoryDescription = "破産寸前の会社を過去最高利益の黒字まで押し上げた実績のあるCTOのブログ　システムの問い合わせの受付やシステム設計の思想・ベストプラクティス　完全に無料で使えるツールの紹介もしています"
	view.CategoryTxt = template.HTML("破産寸前の会社を過去最高利益の黒字まで押し上げた実績のあるCTOのブログ<br>システムの問い合わせの受付やシステム設計の思想・ベストプラクティス<br>完全に無料で使えるツールの紹介もしています")

	view.Note = notes
	view.CategoryList = categoryList
	tpl := template.Must(template.ParseFiles("tpl/category.html"))
	tpl.Execute(w, view)
}
