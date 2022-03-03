package app

import (
	"database/sql"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"../common"
	_ "github.com/lib/pq" // this driver for postgres

	"github.com/grokify/html-strip-tags-go"
)

// Category category page list question
func Category(w http.ResponseWriter, r *http.Request) {

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

	u := strings.Split(r.URL.Path, "/")

	if u[2] != "1" && u[2] != "2" && u[2] != "3" && u[2] != "4" && u[2] != "5" &&
		u[2] != "6" && u[2] != "7" && u[2] != "8" {
		return
	}
	rows, err := db.Query("SELECT * FROM m_category_tree WHERE level_"+u[2]+" = $1", u[3])
	if err != nil {
		log.Print(err) // it is sanitized 
	}
	treeList := map[int]map[string]string{}

	type MCategoryTree struct {
		LeafID    int       // leaf_id
		Level1    int       // level_1
		Level2    int       // level_2
		Level3    int       // level_3
		Level4    int       // level_4
		Level5    int       // level_5
		Level6    int       // level_6
		Level7    int       // level_7
		Level8    int       // level_8
		UpdatedAt time.Time // updated_at
	}
	leaf := false
	d := MCategoryTree{}
	for rows.Next() {
		if err := rows.Scan(&d.LeafID, &d.Level1, &d.Level2, &d.Level3, &d.Level4, &d.Level5, &d.Level6, &d.Level7, &d.Level8, &d.UpdatedAt); err != nil {
			log.Print(err)
		}
		list := map[string]string{}
		switch u[2] {
		case "1":
			list["level"] = "2"
			list["category_name"] = ""
			treeList[d.Level2] = list
		case "2":
			list["level"] = "3"
			list["category_name"] = ""
			treeList[d.Level3] = list
		case "3":
			list["level"] = "4"
			list["category_name"] = ""
			treeList[d.Level4] = list
		case "4":
			list["level"] = "5"
			list["category_name"] = ""
			treeList[d.Level5] = list
		case "5":
			list["level"] = "6"
			list["category_name"] = ""
			treeList[d.Level6] = list
		}
		if strconv.Itoa(d.LeafID) == u[3] {
			leaf = true
		}
	}
	delete(treeList, 0)
	whereIn := u[3]
	forBreadCrumb := map[int]map[string]string{}
	list := map[string]string{}
	if d.Level1 > 0 && u[2] > "1" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level1)
		list = map[string]string{}
		list["level"] = "1"
		list["category_name"] = ""
		forBreadCrumb[d.Level1] = list
	}
	if d.Level2 > 0 && u[2] > "2" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level2)
		list = map[string]string{}
		list["level"] = "2"
		list["category_name"] = ""
		forBreadCrumb[d.Level2] = list
	}
	if d.Level3 > 0 && u[2] > "3" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level3)
		list = map[string]string{}
		list["level"] = "3"
		list["category_name"] = ""
		forBreadCrumb[d.Level3] = list
	}
	if d.Level4 > 0 && u[2] > "4" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level4)
		list = map[string]string{}
		list["level"] = "4"
		list["category_name"] = ""
		forBreadCrumb[d.Level4] = list
	}
	if d.Level5 > 0 && u[2] > "5" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level5)
		list = map[string]string{}
		list["level"] = "5"
		list["category_name"] = ""
		forBreadCrumb[d.Level5] = list
	}
	if d.Level6 > 0 && u[2] > "6" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level6)
		list["level"] = "6"
		list["category_name"] = ""
		forBreadCrumb[d.Level6] = list
	}
	if d.Level7 > 0 && u[2] > "7" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level7)
		list["level"] = "7"
		list["category_name"] = ""
		forBreadCrumb[d.Level7] = list
	}
	if d.Level8 > 0 && u[2] > "8" {
		whereIn = whereIn + "," + strconv.Itoa(d.Level8)
		list["level"] = "8"
		list["category_name"] = ""
		forBreadCrumb[d.Level8] = list
	}
	whereIn2 := u[3] // it is sanitized 
	for i := range treeList {
		whereIn2 = whereIn2 + "," + strconv.Itoa(i)
	}
	if leaf {
		rows, err = db.Query(`SELECT category_id, category_name, category_description FROM m_category_name WHERE category_id = ` +
			u[3] + ` OR category_id in ( SELECT level_1 FROM m_category_tree GROUP BY level_1 )`)
	} else {
		rows, err = db.Query("SELECT category_id, category_name, category_description FROM m_category_name WHERE category_id in (" + whereIn + "," + whereIn2 + ")")
	}

	if err != nil {
		log.Print(err)
	}
	type MCategoryName struct {
		CategoryID          int       // category_id
		CategoryName        string    // category_name
		UpdatedAt           time.Time // updated_at
		CategoryDescription string    // category_description
	}
	var categoryList []CategoryList
	for rows.Next() {
		r := MCategoryName{}
		if err := rows.Scan(&r.CategoryID, &r.CategoryName, &r.CategoryDescription); err != nil {
			log.Print(err)
		}
		if _, ok := forBreadCrumb[r.CategoryID]; ok {
			forBreadCrumb[r.CategoryID]["category_name"] = r.CategoryName
		}
		if _, ok := treeList[r.CategoryID]; ok {
			treeList[r.CategoryID]["category_name"] = r.CategoryName
		}
		if u[3] == strconv.Itoa(r.CategoryID) {
			view.CategoryName = r.CategoryName
			view.CategoryDescription = r.CategoryDescription
			view.CategoryTxt = template.HTML(strings.Replace(r.CategoryDescription, "\n", "<br>", -1))
		} else if leaf {
			y := CategoryList{}
			y.Level = "1"
			y.CategoryID = r.CategoryID
			y.CategoryName = r.CategoryName
			categoryList = append(categoryList, y)
		}
	}

	var breadCrumb []BreadCrumb
	for i, v := range forBreadCrumb {
		y := BreadCrumb{}
		y.Level = v["level"]
		y.CategoryID = i
		y.CategoryName = v["category_name"]
		breadCrumb = append(breadCrumb, y)
	}
	sort.Slice(breadCrumb, func(i, j int) bool { return breadCrumb[i].Level < breadCrumb[j].Level }) // DESC
	view.BreadCrumb = breadCrumb

	if !(leaf) {

		// var categoryList []CategoryList
		for i, v := range treeList {
			y := CategoryList{}
			y.Level = v["level"]
			y.CategoryID = i
			y.CategoryName = v["category_name"]
			categoryList = append(categoryList, y)
		}
	}
// fmt.Printf("categoryList %#v\n", categoryList)
	if leaf {
		rows, err = db.Query(`SELECT note_id, note_title, note_txt, updated_at 
			FROM t_note WHERE category_id = ` + u[3] + `ORDER BY note_id DESC`)
	} else {
		rows, err = db.Query(`SELECT note_id, note_title, note_txt, updated_at
			FROM t_note WHERE list_category_id = ` + u[3] + `ORDER BY note_id DESC`)
	}

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
		r.UpdatedAt = ti.Format("2006年1月2日")
		r.NoteTxt = strip.StripTags(txt)[1:256]
		notes = append(notes, r)
	}
	
	view.Note = notes
	view.CategoryList = categoryList
	tpl := template.Must(template.ParseFiles("tpl/category.html"))
	tpl.Execute(w, view)
}
