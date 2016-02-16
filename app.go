package main
import (
  "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/drone/routes"
  "time"
  "encoding/json"
)

type Page struct {
  title string
  body []byte
}

type Rental struct {
    City string
    Owner string
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

func checkErr(err error) {
    if err != nil {
      log.Fatal(err)
    }
}

func handler(w http.ResponseWriter, r *http.Request)  {
 // params := r.URL.Query()
  //uid := params.Get(":uid")
  //var owner string
  rows, err := db.Query("SELECT * FROM rentals;")
  if err == sql.ErrNoRows {
    log.Fatal("No Results Found")
  }
  if err != nil {
    log.Fatal(err)
  }
  rentals := []Rental{}

  for rows.Next() {
      var owner string
      var id int
      var city string
      var bedrooms int
      var createdAt time.Time
      var updatedAt time.Time
      err = rows.Scan(&id, &city, &owner, &bedrooms, &createdAt, &updatedAt)
      checkErr(err)
      rentals = append(rentals, Rental{City: city, Owner: owner})
      //fmt.Fprintf(w, "<h1>%s</h1><p>Onwer: %s</p>", city, owner)
  }
  js, err := json.Marshal(rentals)
  fmt.Fprintf(w, "<p>%s</p>", js)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, _ := loadPage(title)
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)
}

func main() {
  mux := routes.New()
  mux.Get("/", handler)
  //http.HandleFunc("/view/", viewHandler)
  //http.HandleFunc("/", handler)
  http.Handle("/", mux)
  http.ListenAndServe(":3000", nil)
}
