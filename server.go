package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.gohtml"))
}

func server() {
	fmt.Println("running")
	http.HandleFunc("/", getHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalln(err)
		return
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	var files []string
	// file, err := os.Open("dst")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer file.Close()

	if err := filepath.Walk("dst", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, strings.Split(path, "/")[1])
		return nil
	}); err != nil {
		panic(err)
	}

	tpl.ExecuteTemplate(w, "index.gohtml", files)
}
