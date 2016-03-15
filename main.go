package main

  import (
      "fmt"
      "net/http"
      "html/template"
      "encoding/json"
      "strings"
      "regexp"
      "log"
  )

  type Asset struct{
      Title string
      Subtitle string
  }

  /*Return a Single Asset according a title*/
  func  GetAsset(title string) (*Asset, error){
    //Mock result of SearchAWS of Valerio's method
    const searchAws =
    `
      {"Title":"title1", "Subtitle": "example subtitle 1" }
      {"Title":"title2", "Subtitle": "example subtitle 2" }
    `
    var a Asset
    dec := json.NewDecoder(strings.NewReader(searchAws))
    err := dec.Decode(&a);
    if err != nil{
        log.Fatal(err)
    }
    return &a, nil
  }

  func editHandler(w http.ResponseWriter, r *http.Request, title string){
      p, err := GetAsset(title)
      if( err != nil){
        fmt.Printf("%s\n", err)
      }
      renderTemplate(w,"edit",p)
  }

  func viewHandler(w http.ResponseWriter, r *http.Request, title string){
      p, err := GetAsset(title)
      if( err != nil){
         http.Redirect(w,r,"edit/"+title,http.StatusFound)
         return
      }
      renderTemplate(w,"view",p)
  }


  var templates = template.Must(template.ParseFiles("edit.html","view.html"))

  func renderTemplate(w http.ResponseWriter, templateName string, a *Asset){
      err := templates.ExecuteTemplate(w, templateName +".html", a)
      if err!= nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
      }
  }

  /*restrict access to edit and view page with alphanumeric query*/
  var validPath = regexp.MustCompile("(edit|view)/([a-zA-Z0-9]+)$")

  func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
      return func(w http.ResponseWriter, r *http.Request){
          m := validPath.FindStringSubmatch(r.URL.Path)
          if m == nil {
              fmt.Printf(r.URL.Path + " not found\n")
              http.NotFound(w,r)
              return
          }
          fn(w,r,m[2])
      }
  }

  func main(){
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.ListenAndServe(":8080",nil)
  }
