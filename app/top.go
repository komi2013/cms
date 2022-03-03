package app

import (
	"database/sql"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	rows, err := db.Query("SELECT level_1 FROM m_category_tree GROUP BY level_1")
	if err != nil {
		log.Print(err)
	}
	treeList := map[int]map[string]string{}
	whereIn := "0"
	for rows.Next() {
		level_1 := 0
		if err := rows.Scan(&level_1); err != nil {
			log.Print(err)
		}
		list := map[string]string{}
		list["level"] = "1"
		list["category_name"] = ""
		treeList[level_1] = list
		whereIn = whereIn + "," + strconv.Itoa(level_1)
	}
	delete(treeList, 0)

	rows, err = db.Query("SELECT category_id, category_name FROM m_category_name WHERE category_id in (" + whereIn + ")")
	if err != nil {
		log.Print(err)
	}
	type MCategoryName struct {
		CategoryID          int       // category_id
		CategoryName        string    // category_name
	}
	for rows.Next() {
		r := MCategoryName{}
		if err := rows.Scan(&r.CategoryID, &r.CategoryName); err != nil {
			log.Print(err)
		}
		if _, ok := treeList[r.CategoryID]; ok {
			treeList[r.CategoryID]["category_name"] = r.CategoryName
		}
	}

	var categoryList []CategoryList
	for i, v := range treeList {
		y := CategoryList{}
		y.Level = v["level"]
		y.CategoryID = i
		y.CategoryName = v["category_name"]
		categoryList = append(categoryList, y)
	}
	rows, err = db.Query(`SELECT note_id, note_title, note_txt, updated_at
			FROM t_note WHERE list_category_id = 0 ORDER BY note_id DESC`)
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
		r.NoteTxt = strip.StripTags(txt)[1:256]
		notes = append(notes, r)
	}
	view.CategoryName = "炎上案件上等CTO"
	view.CategoryDescription = "炎上案件上等CTOのブログ　システムの問い合わせの受付やシステム設計の思想・ベストプラクティス　完全に無料で使えるツールの紹介もしています"
	view.CategoryTxt = template.HTML("炎上案件上等CTOのブログ<br>システムの問い合わせの受付やシステム設計の思想・ベストプラクティス<br>完全に無料で使えるツールの紹介もしています")

	view.Note = notes
	view.CategoryList = categoryList
	tpl := template.Must(template.ParseFiles("tpl/category.html"))
	tpl.Execute(w, view)
}
