package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		panic(err.Error)
	}
}

type filestruct struct {
	Name string
	Data []byte
}

type templdata struct {
	Callback string
	Files    string
}

// Serves Files from the "upload" folder as JSON array
// via JSONP
func GetFiles(w http.ResponseWriter, r *http.Request) {
	// TODO: Read files from upload folder
	w.Header().Set("Access-Control-Allow-Origin", "*")

	files_info, err := ioutil.ReadDir("./upload")
	CheckError(err)

	files := make([]filestruct, len(files_info))

	var data []byte

	for i, f := range files_info {
		if f.IsDir() {
			continue
		}
		data, err = ioutil.ReadFile("./upload/" + f.Name())
		CheckError(err)
		files[i] = filestruct{Name: f.Name(), Data: data}
		os.Remove("./upload/" + f.Name())
	}

	tmpl, err := template.ParseFiles("getfiles.jsonp")
	CheckError(err)

	files_json, err := json.Marshal(files)
	CheckError(err)

	err = tmpl.Execute(w, templdata{Callback: mux.Vars(r)["callback"], Files: string(files_json)})
	CheckError(err)
}

// Serves the file index.html
func Index(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("index.html")
	defer file.Close()
	CheckError(err)

	io.Copy(w, file)
}

func ServeJS(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(mux.Vars(r)["js_filename"] + ".js")
	defer file.Close()
	CheckError(err)

	io.Copy(w, file)
}

// Recieves Files via Multipart upload and saves them in the "download" Folder
func Upload(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		files []filestruct
	)

	err = json.NewDecoder(r.Body).Decode(&files)
	CheckError(err)

	for _, f := range files {
		err = ioutil.WriteFile("./download/"+f.Name, f.Data, 0666)
		CheckError(err)
	}

	fmt.Fprintf(w, "%s", "Got it!")
}

func main() {

	mx := mux.NewRouter()
	mx.HandleFunc("/{js_filename}.js", ServeJS)
	mx.HandleFunc("/getfiles/callback={callback}", GetFiles)
	mx.HandleFunc("/index", Index)
	mx.HandleFunc("/upload", Upload)

	n := negroni.Classic()
	n.UseHandler(cors.Default().Handler(mx))

	n.Run(":" + os.Getenv("PORT"))
}
