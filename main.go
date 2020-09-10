// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
	"log"

)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "posts/"+p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "posts/"+title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}



func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}




var templates = template.Must(template.ParseFiles("tmp/edit.html", "tmp/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Fprintf(w,"hi")
		return;
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view|listing)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			
			fmt.Fprintf(w,"hi")

			return
		}
		fn(w, r, m[2])
	}
}

func viewListing(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"wag1")
}

type Listing struct{
	Title string
	Listings []string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tem, err:= template.ParseFiles("tmp/frontpage.html")
	if err !=nil{
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
		return
	}

	//lets get all .txt files in posts as they are posts

	posts,err := ioutil.ReadDir("./posts")
	if err!=nil{
		log.Fatal(err)
	}

	var allPosts []string


	for _, post:= range posts{
		allPosts=append(allPosts,post.Name())
	}








	data:=Listing{
		Title:"Listing page",
		Listings:allPosts,
	}
	tem.Execute(w,data)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":80", nil)
}