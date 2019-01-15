package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bukodi/go-playgroud/swaggerui"
	"github.com/emicklei/go-restful"
	openapi "github.com/emicklei/go-restful-openapi"
	// "github.com/emicklei/go-restful-swagger12"
)

// Book is an example entity
type Book struct {
	Title  string
	Author string
}

var books []Book

//var book *Book

func main() {

	books = []Book{
		Book{"Egri csillagok", "Gárdonyi Géza"},
		Book{"Arany ember", "Jókai Mór"},
	}

	container := restful.NewContainer()

	ws := new(restful.WebService)
	ws.Path("/api/books")
	ws.Doc("Minden cica aranyos")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(getBooks).
		Doc("Search all books").
		Param(ws.QueryParameter("language", "en,nl,de").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "last known timestamp").DataType("datetime")))

	ws.Route(ws.GET("/{medium}").To(noop).
		Doc("Search all books").
		Param(ws.PathParameter("medium", "digital or paperback").DataType("string")).
		Param(ws.QueryParameter("language", "en,nl,de").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "last known timestamp").DataType("datetime")).
		Do(returns200, returns500))

	ws.Route(ws.PUT("/{medium}").To(noop).
		Doc("Add a new book").
		Param(ws.PathParameter("medium", "digital or paperback").DataType("string")).
		Reads(Book{}))

	container.Add(ws)

	fmt.Printf("%#v", ws)

	// You can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := openapi.Config{
		WebServices: container.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json",
	}
	container.Add(openapi.NewOpenAPIService(config))

	log.Print("start listening on localhost:8080")

	http.HandleFunc("/swagger-ui/", swaggeruiHandler)
	//	http.HandleFunc("/charts", func(w http.ResponseWriter, r *http.Request) {
	//		http.ServeFile(w, r, "/index.html")
	//	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/charts" || r.URL.Path == "/forms" {
			serveIndexHTML(w, r)
			return
		} else {
			http.FileServer(assetFS()).ServeHTTP(w, r)
		}
	})
	//http.Handle("/")

	server := &http.Server{Addr: ":8080", Handler: container}
	log.Fatal(server.ListenAndServe())
}

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

func noop(req *restful.Request, resp *restful.Response) {}

func getBooks(req *restful.Request, resp *restful.Response) {
	resp.WriteEntity(books)
}

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", books)
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	f, err := assetFS().Open("/index.html")
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	// serveContent will check modification time
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func swaggeruiHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/swagger-ui/") {
		http.NotFound(w, r)
		return
	}

	assetName := r.URL.Path[len("/swagger-ui/"):]
	bytes, err := swaggerui.Asset(assetName)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(assetName, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}

	w.Write(bytes)
	//r.URL.Path
}
