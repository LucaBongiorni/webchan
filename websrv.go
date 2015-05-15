package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
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
	files_info, err := ioutil.ReadDir("./upload")
	CheckError(err)

	files := make([]filestruct, len(files_info))

	var data []byte

	for i, f := range files_info {
		data, err = ioutil.ReadFile("./upload/" + f.Name())
		CheckError(err)
		files[i] = filestruct{Name: f.Name(), Data: data}
		os.Remove("./upload/" + f.Name())
	}

	//filesstruct := []filestruct{{Name: "onefile", Data: []byte("Super Inhalt der in \"onefile\" steht.")}, {Name: "anotherfile.txt", Data: []byte("Noch mehr cooler Scheisz!")}}

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
	r.ParseMultipartForm(20000)
	var (
		file *os.File
		err  error
		data []byte
	)

	for key, value := range r.MultipartForm.Value {
		file, err = os.OpenFile("./download/"+key, os.O_WRONLY|os.O_CREATE, 06600)
		CheckError(err)

		data, err = base64.StdEncoding.DecodeString(value[0])
		CheckError(err)

		_, err = file.Write(data)
		CheckError(err)

		err = file.Close()
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

	http.ListenAndServe(":"+os.Getenv("PORT"), mx)
}
