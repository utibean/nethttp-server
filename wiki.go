package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

//
//func main() {
//	p1 := Page{Title: "go-in-action", Body: []byte("This is a simple page")}
//	p1.save()
//	p2, _ := loadPage(p1.Title)
//	fmt.Println(string(p2.Body))
//}

func main() {
	//http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTitle(writer http.ResponseWriter, request *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(request.URL.Path)
	if m == nil {
		http.NotFound(writer, request)
		return "", errors.New("invalid page title")
	}
	return m[2], nil
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		m := validPath.FindStringSubmatch(request.URL.Path)
		if m == nil {
			http.NotFound(writer, request)
			return
		}
		fn(writer, request, m[2])
	}
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderHtml(writer, "edit", page)
}

func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderHtml(writer, "view", p)
}

func renderHtml(writer http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(writer, tmpl+".html", p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

//func renderHtml(writer http.ResponseWriter, tmpl string, p *Page) {
//	t, err := template.ParseFiles(tmpl + ".html")
//	if err != nil {
//		http.Error(writer, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = t.Execute(writer, p)
//	if err != nil {
//		http.Error(writer, err.Error(), http.StatusInternalServerError)
//	}
//}

//func viewHandler(writer http.ResponseWriter, request *http.Request) {
//	title := request.URL.Path[len("/view/"):]
//	p, _ := loadPage(title)
//	fmt.Fprintf(writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
//}

//
//func handler(writer http.ResponseWriter, request *http.Request) {
//	fmt.Fprintf(writer, "Hi, I love %s!", request.URL.Path[1:])
//}
