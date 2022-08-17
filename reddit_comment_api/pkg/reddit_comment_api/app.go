package reddit_comment_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	db *sqlx.DB
}

type Comment struct {
	ID         string
	parent_ID  string
	post_ID    string
	body       string
	author     string
	updated_at int
}

type SubmissionNode struct {
	ID         string
	Body       string
	Author     string
	Updated_at int
	Children   []*SubmissionNode
}

func NewApp(db *sqlx.DB) App {
	a := App{
		db: db,
	}
	return a
}

// REDDIT COMMENT PREFIXES:
// #type prefixes
// # t1_	Comment
// # t2_	Account
// # t3_	Link
// # t4_	Message
// # t5_	Subreddit
// # t6_	Award

func (a *App) GetLatestSubmissions(w http.ResponseWriter, r *http.Request) {
	var results []SubmissionNode
	queryStatement := "SELECT id, body, author, updated_at FROM submission ORDER BY updated_at DESC LIMIT 10;"
	rows, err := a.db.Queryx(queryStatement)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var result SubmissionNode
		err := rows.StructScan(&result)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) GetLatestComments(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	post_ID := path[2]

	var results SubmissionNode
	queryStatement := "SELECT id, body, author, updated_at FROM submission WHERE ID = $1 LIMIT 1;"
	if err := a.db.QueryRowx(queryStatement, post_ID).StructScan(&results); err != nil {
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err := a.appendChildren(&results)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) appendChildren(parent *SubmissionNode) error {
	var results []*SubmissionNode
	queryStatement := "SELECT id, body, author, created_at, updated_at FROM comments WHERE PARENT_ID = $1 ORDER BY updated_at desc;"
	if err := a.db.Select(&results, queryStatement, parent.ID); err != nil {
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	parent.Children = append(parent.Children, results...)

	for _, v := range parent.Children {
		if v.ID == "t1_ikjzuir" {
			fmt.Println(v)
		}
		a.appendChildren(v)
		fmt.Println(parent)
	}

	return nil
}

func (a *App) Start() error {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/submissions/", a.GetLatestSubmissions).Methods("GET")
	rtr.HandleFunc("/submissions/{^.*[a-zA-Z0-9]+.*$}/comments", a.GetLatestComments).Methods("GET")

	err := http.ListenAndServe(":8080", rtr)
	return err
}
