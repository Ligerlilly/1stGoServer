package main
import (
  "fmt"
  "net/http"
  "io/ioutil"
)

type Page struct {
  title string
  body []byte
}

func (p *Page) save() error {
  filename := p.title + ".txt"
  return ioutil.WriteFile(filename, p.body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{title: title, body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintf(w, "Hello world")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, _ := loadPage(title)
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/", handler)
  http.ListenAndServe(":3000", nil)
}
