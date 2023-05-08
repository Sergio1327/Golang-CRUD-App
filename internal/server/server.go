package server

import (
	"CRUD/internal/db"
	"CRUD/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Run() {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	router := http.NewServeMux()
	router.HandleFunc("/get", GetAllCV(db))
	router.HandleFunc("/get/", GetById(db))
	router.HandleFunc("/create", Create(db))
	router.HandleFunc("/update/", Update(db))
	router.HandleFunc("/delete/", Delete(db))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetAllCV(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			rows, err := db.Query("select * from jobs")
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Fprint(w, "no data in database")
				}
				log.Fatal(err)
			}
			defer rows.Close()

			var CVs []model.CV

			for rows.Next() {
				var c model.CV
				err := rows.Scan(&c.Name, &c.Job, &c.Salary)
				if err != nil {
					log.Fatal(err)
				}
				CVs = append(CVs, c)
			}
			json.NewEncoder(w).Encode(CVs)
		} else {
			http.Error(w,"error,enter correct method",http.StatusMethodNotAllowed)
		}
	}
}

func Create(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var user model.CV
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec("insert into jobs(name,job,salary) values ($1,$2,$3)",
				user.Name, user.Job, user.Salary)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, "data was succesful saved in database")
		} else {
			http.Error(w,"error,eneter correct method",http.StatusMethodNotAllowed)
		}
	}
}

func GetById(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := strings.TrimPrefix(r.URL.Path, "/get/")
			row, err := db.QueryRow("select * from jobs where name=$1", id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Fprint(w, "error, no user with this id")
				}
				log.Fatal(err)
			}
			var u model.CV
			if err := row.Scan(&u.Name, &u.Job, &u.Salary); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "no data in database with this id", http.StatusNotFound)
				}
			}
			json.NewEncoder(w).Encode(u)
		} else {
			http.Error(w, "error,please enter correct method", http.StatusMethodNotAllowed)
		}
	}
}
func Update(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			name := strings.TrimPrefix(r.URL.Path, "/update/")
			var u model.CV
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				http.Error(w, "error", http.StatusBadRequest)
			}
			_, err := db.Exec("update jobs set name=$1,job=$2,salary=$3 where name=$4", u.Name, u.Job, u.Salary, name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, "succesfully updated")
		} else {
			http.Error(w, "error,please select correct method", http.StatusMethodNotAllowed)
		}
	}
}

func Delete(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			name := strings.TrimPrefix(r.URL.Path, "/delete/")
			_, err := db.Exec("delete from jobs where name=$1", name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, "succesfully deleted")
		} else {
			http.Error(w, "not correct method", http.StatusMethodNotAllowed)
		}
	}
}
