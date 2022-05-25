package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"../common"
	_ "github.com/lib/pq" // this driver for postgres
)

func Note(w http.ResponseWriter, r *http.Request) {

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
		UpdatedAt     time.Time
		NoteTitle     string
		NoteTxt       template.HTML
		CategoryID    string
		NoteImg       string
		WroteDate     string
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
		Note                Note
		URL                 string
	}
	var view View
	view.CacheV = common.CacheV
	view.CSRF = ""

	u := strings.Split(r.URL.Path, "/")

	stmt, err := db.Prepare(`SELECT note_id, note_title, note_txt, updated_at, category_id, note_img 
			FROM t_note WHERE note_id = $1`)
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()
	note := Note{}
	err = stmt.QueryRow(u[2]).Scan(&note.NoteID, &note.NoteTitle, &note.NoteTxt, &note.UpdatedAt, &note.CategoryID, &note.NoteImg)
	note.WroteDate = note.UpdatedAt.Format("2006年01月02日")
	if err != nil {
		log.Print(err)
	}
	// fmt.Printf("note %#v\n", note)
	view.Note = note
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
	convert := make(map[int]string)
	var tree [][2]int
	var x [2]int
	rows, err := db.Query("SELECT * FROM m_category_tree WHERE leaf_id = $1", note.CategoryID)
	if err != nil {
		log.Print(err)
	}
	whereIn := ""
	for rows.Next() {
		r := MCategoryTree{}
		if err := rows.Scan(&r.LeafID, &r.Level1, &r.Level2, &r.Level3, &r.Level4, &r.Level5, &r.Level6, &r.Level7, &r.Level8, &r.UpdatedAt); err != nil {
			log.Print(err)
		}
		whereIn = strconv.Itoa(r.Level1)
		x[0] = 1
		x[1] = r.Level1
		tree = append(tree, x)
		convert[r.Level1] = ""
		if r.Level2 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level2)
			x[0] = 2
			x[1] = r.Level2
			tree = append(tree, x)
			convert[r.Level2] = ""
		}
		if r.Level3 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level3)
			x[0] = 3
			x[1] = r.Level3
			tree = append(tree, x)
			convert[r.Level3] = ""
		}
		if r.Level4 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level4)
			x[0] = 4
			x[1] = r.Level4
			tree = append(tree, x)
			convert[r.Level4] = ""
		}
		if r.Level5 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level5)
			x[0] = 5
			x[1] = r.Level5
			tree = append(tree, x)
			convert[r.Level5] = ""
		}
		if r.Level6 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level6)
			x[0] = 6
			x[1] = r.Level6
			tree = append(tree, x)
			convert[r.Level6] = ""
		}
		if r.Level7 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level7)
			x[0] = 7
			x[1] = r.Level7
			tree = append(tree, x)
			convert[r.Level7] = ""
		}
		if r.Level8 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level8)
			x[0] = 8
			x[1] = r.Level8
			tree = append(tree, x)
			convert[r.Level8] = ""
		}
	}
	type MCategoryName struct {
		CategoryID          int       // category_id
		CategoryName        string    // category_name
		UpdatedAt           time.Time // updated_at
		CategoryDescription string    // category_description
	}
	if whereIn != "" {
		rows, err = db.Query("SELECT category_id, category_name FROM m_category_name WHERE category_id in (" + whereIn + ")")
		if err != nil {
			log.Print(err)
		}
		for rows.Next() {
			r := MCategoryName{}
			if err := rows.Scan(&r.CategoryID, &r.CategoryName); err != nil {
				log.Print(err)
			}
			convert[r.CategoryID] = r.CategoryName
		}
		var breadCrumb []BreadCrumb
		for _, v := range tree {
			y := BreadCrumb{}
			y.Level = strconv.Itoa(v[0])
			y.CategoryID = v[1]
			y.CategoryName = convert[v[1]]
			breadCrumb = append(breadCrumb, y)
		}

		sort.Slice(breadCrumb, func(i, j int) bool { return breadCrumb[i].Level < breadCrumb[j].Level }) // DESC
		// fmt.Printf("breadCrumb %#v\n", breadCrumb)
		view.BreadCrumb = breadCrumb
	}
	rows, err = db.Query(`SELECT category_id, category_name FROM m_category_name WHERE category_id in 
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
	view.CategoryList = categoryList
	view.URL = r.Host + r.URL.Path
	fmt.Println(r.Host)
	log.Print("mukashi")
	tpl := template.Must(template.ParseFiles("tpl/note.html"))
	tpl.Execute(w, view)

}
