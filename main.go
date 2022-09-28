package main

import (
	"context"
	"fmt"
	"golang-manipulate/connection"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// struct
type Struktur struct {
	Name string
	Start_date string
	End_date string
	Deskripsi string
	Node string
	React string
	Laravel string
	Golang string
	Gambar string
	Duration string
	Id int
}

var iniArray = []Struktur{}
func main() {
	route := mux.NewRouter()
	connection.Dbkonek()
	// path prefix
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	// routing

	route.HandleFunc("/",home).Methods("GET")
	route.HandleFunc("/add-blog", addProject).Methods("GET")
	route.HandleFunc("/store-blog", storeProject).Methods("POST")
	route.HandleFunc("/detail-blog/{id}", detailProject).Methods("GET")
	route.HandleFunc("/edit/{id}", editProject).Methods("GET")
	route.HandleFunc("/update-blog/{id}", updateProject).Methods("POST")
	route.HandleFunc("/delete/{id}", deleteProject).Methods("GET")


	// contact Route

	// route.HandleFunc("/contact", contact).Methods("GET")

	// server
	fmt.Println("server is runing on 127.0.0.1:5000")
	http.ListenAndServe("127.0.0.1:5000",route)
}

func home(res http.ResponseWriter, req *http.Request)  {
	res.Header().Set("Content-Type","text/html; charset=utf-8")
	theme, err := template.ParseFiles("views/blog/index.html")

	if err != nil {
		res.Write([]byte("massage : HACKER JANGAN MENYERANG !" + err.Error()))
	}

	data,err := connection.Konekdb.Query(context.Background(), "SELECT id, name, description FROM tb_projects")

	var result []Struktur

	for data.Next(){
		var each = Struktur{}

		err := data.Scan(&each.Id, &each.Name, &each.Deskripsi)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}

		mapping := map[string]interface{}{
			"show" :result,
		}

	theme.Execute(res, mapping)
}

func addProject(res http.ResponseWriter, req *http.Request)  {
	res.Header().Set("Content-Type","text/html; charset=utf-8")
	theme, err := template.ParseFiles("views/blog/addproject.html")

	if err != nil {
		res.Write([]byte("massage : HACKER JANGAN MENYERANG !" + err.Error()))
	}

	theme.Execute(res, nil)
}
func storeProject(res http.ResponseWriter, req *http.Request)  {
	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := req.PostForm.Get("title")
	// start_date := req.PostForm.Get("start-date")
	// end_date := req.PostForm.Get("end-date")
	desc := req.PostForm.Get("desc")
	// node := req.PostForm.Get("node")
	// laravel := req.PostForm.Get("laravel")
	// react := req.PostForm.Get("react")
	// // golang := req.PostForm.Get("golang")

	// layouts := "2006-01-02"
	// convStartDate, _ := time.Parse(layouts, start_date)  
	// convEndtDate, _ := time.Parse(layouts, end_date)  

	// hourse := convEndtDate.Sub(convStartDate).Hours()
	// days := hourse/24
	// weeks := days/7
	// months := days/30
	// years := months/12

	// // var duration string
	// if days >= 1 && days <= 6 {
    //     duration = strconv.Itoa(int(days)) + " days"
    // } else if days >= 7 && days <= 29 {
    //     duration = strconv.Itoa(int(weeks)) + " weeks"
    // } else if days >= 30 && days <= 364 {
    //     duration = strconv.Itoa(int(months)) + " months"
    // } else if days >= 365 {
    //     duration = strconv.Itoa(int(years)) + " years"
    // }

	_, err = connection.Konekdb.Exec(context.Background(), "INSERT INTO tb_projects(name, description) VALUES ($1,$2)",title, desc)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("message : " + err.Error()))
		return
	}


	// var newProject = Struktur{
	// 	Name : title,
	// 	Start_date : start_date,
	// 	End_date :  end_date,
	// 	Deskripsi : desc,
	// 	Node : node,
	// 	Laravel : laravel,
	// 	React : react,
	// 	Golang : golang,
	// 	Duration : duration,
	// 	Id : len(iniArray),
	// }

	// iniArray = append(iniArray, newProject)
	// fmt.Println(iniArray)

	http.Redirect(res, req, "/", http.StatusMovedPermanently)
}



// masih eror di ID NYAAA
func detailProject(res http.ResponseWriter, req *http.Request)  {

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	theme, err := template.ParseFiles("views/blog/detail.html")

	if err != nil {
		res.Write([]byte("Hacker jangan menyerang! :" + err.Error()))
		return
	}
	var blogDetail = Struktur{}

	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	err = connection.Konekdb.QueryRow(context.Background(), " SELECT id, name, description FROM tb_projects WHERE id=$1", id).Scan(&blogDetail.Id, &blogDetail.Name, &blogDetail.Deskripsi)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte ("message ; " + err.Error()))
	}

	data := map[string]interface{}{
		"detail" : blogDetail,
	}
	fmt.Println(data)
	theme.Execute(res, data)
}
func editProject(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/blog/edit-project.html")

	if err != nil {
		res.Write([]byte("message : "+ err.Error()))
		return
	}

	var editProject = Struktur{}

	index, _ := strconv.Atoi(mux.Vars(req)["id"])

	for i, project := range iniArray {
		if index == i {
			editProject = Struktur{
				Name: project.Name,
				Deskripsi: project.Deskripsi,
				Start_date: project.Start_date,
				End_date: project.End_date,
				Node: project.Node,
				Golang: project.Golang,
				React: project.React,
				Laravel: project.Laravel,
				Id: project.Id,
			}
		}

	}

	data := map[string]interface{}{
		"EditProject": editProject,
	}

	tmpl.Execute(res, data)
}

func updateProject(res http.ResponseWriter, req *http.Request){
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	
	err := req.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := req.PostForm.Get("nameProject")
	description := req.PostForm.Get("description")
		EditdataProject := Struktur{
		Name: title,
		Deskripsi: description,
		Id: id,
	}

	iniArray[id] = EditdataProject

	http.Redirect(res, req, "/", http.StatusFound)
}

func deleteProject(res http.ResponseWriter, req *http.Request){
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	_,err := connection.Konekdb.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1",id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(res, req, "/", http.StatusMovedPermanently)
}