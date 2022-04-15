package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//
//func main() {
//	p1 := Page{Title: "go-in-action", Body: []byte("This is a simple page")}
//	p1.save()
//	p2, _ := loadPage(p1.Title)
//	fmt.Println(string(p2.Body))
//}

func main() {
	//http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/save/"):]
	body := request.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderHtml(writer, "edit", page)
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderHtml(writer, "view", p)
}

func renderHtml(writer http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(writer, p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

//func viewHandler(writer http.ResponseWriter, request *http.Request) {
//	title := request.URL.Path[len("/view/"):]
//	p, _ := loadPage(title)
//	fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
//}

//
//func handler(writer http.ResponseWriter, request *http.Request) {
//	fmt.Fprintf(writer, "Hi, I love %s!", request.URL.Path[1:])
//}
