package main
import (
  "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

type Page struct {
  title string
  body []byte
}

var db *sql.DB

func init() {
  var err error
  db, err = sql.Open("postgres", "postgres://username:password@localhost/rentals?sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }
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
  var owner string
  err := db.QueryRow("SELECT owner FROM rentals WHERE city = 'Detroit'").Scan(&owner)
  if err == sql.ErrNoRows {
    log.Fatal("No Results Found")
  }
  if err != nil {
    log.Fatal(err)
  }
  fmt.Fprintf(w, "<h1>%s</h1>",owner)
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
