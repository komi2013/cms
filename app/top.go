package app

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../common"
	_ "github.com/lib/pq" // this driver for postgres
)

// Top page almost same with cateogory
func Top(w http.ResponseWriter, r *http.Request) {

	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	type NewNotes struct {
		QuestionID    string
		QuestionTitle string
	}
	type CategoryList struct {
		Level            string
		CategoryID       int
		CategoryName     string
		NewNotes []NewNotes
	}
	type List struct {
		Level1 int
	}
	type View struct {
		CacheV              string
		CSRF                string
		CategoryList        []CategoryList
		CategoryName        string
		CategoryDescription string
		CategoryTxt         template.HTML
		BreadCrumb          []int
		NewNotes    []NewNotes
	}
	var view View
	view.CacheV = common.CacheV
	view.CSRF = common.MakeCSRF(w, r)

	if r.URL.Path != "/" {
		log.Print(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	rows, err := db.Query("SELECT level_1 FROM m_category_tree GROUP BY level_1")
	if err != nil {
		log.Print(err)
	}
	treeList := map[int]map[string]string{}
	d := List{}
	for rows.Next() {
		if err := rows.Scan(&d.Level1); err != nil {
			log.Print(err)
		}
		list := map[string]string{}
		list["level"] = "1"
		list["category_name"] = ""
		treeList[d.Level1] = list
	}
	whereIn2 := "0"
	for i := range treeList {
		whereIn2 = whereIn2 + "," + strconv.Itoa(i)
	}
	rows, err = db.Query("SELECT category_id, category_name, category_description FROM m_category_name WHERE category_id in (" + whereIn2 + ")")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		r := common.MCategoryName{}
		if err := rows.Scan(&r.CategoryID, &r.CategoryName, &r.CategoryDescription); err != nil {
			log.Print(err)
		}
		if _, ok := treeList[r.CategoryID]; ok {
			treeList[r.CategoryID]["category_name"] = r.CategoryName
		}
		if strconv.Itoa(r.CategoryID) == "0" {
			view.CategoryName = r.CategoryName
			view.CategoryDescription = r.CategoryDescription
			view.CategoryTxt = template.HTML(strings.Replace(r.CategoryDescription, "\n", "<br>", -1))
		}
	}
	var newNotesPre []common.MNewNotes
	rows, err = db.Query("SELECT question_id, category_id, question_title, in_list FROM m_category_question WHERE category_id in (" + whereIn2 + ")")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		r := common.MNewNotes{}
		if err := rows.Scan(&r.QuestionID, &r.CategoryID, &r.QuestionTitle, &r.InList); err != nil {
			log.Print(err)
		}
		newNotesPre = append(newNotesPre, r)
	}

	var categoryList []CategoryList
	for i, v := range treeList {
		y := CategoryList{}
		y.Level = v["level"]
		y.CategoryID = i
		y.CategoryName = v["category_name"]
		var newNotes []NewNotes
		for _, v2 := range newNotesPre {
			if y.CategoryID == v2.CategoryID && v2.InList == 1 {
				y2 := NewNotes{}
				y2.QuestionID = strconv.Itoa(v2.QuestionID)
				y2.QuestionTitle = v2.QuestionTitle
				newNotes = append(newNotes, y2)
			}
		}
		y.NewNotes = newNotes
		categoryList = append(categoryList, y)
	}
	view.CategoryList = categoryList
	tpl := template.Must(template.ParseFiles("tpl/category.html"))
	tpl.Execute(w, view)
}
