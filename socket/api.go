//
// @see: https://reinbach.com/golang-webapps-1.html
// 	     https://medium.com/@ivanderbyl/why-you-don-t-need-socket-io-6848f1c871cd

package socket

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"ufuchs/itplus/base/fcc"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

type (
	Context struct {
		Title  string
		Static string
	}
)

//
//
//
func Run(port int, hub *Hub) {

	go hub.run()

	thePort := strconv.Itoa(port)

	http.HandleFunc("/ws", socketService(hub))
	http.HandleFunc("/", home)
	http.HandleFunc(STATIC_URL, staticHandler)

	fmt.Println("Socket: listening on port", thePort)

	if err := http.ListenAndServe(":"+thePort, nil); err != nil {
		fcc.Fatal("Socket ListenAndServe: ", err)
	}

}

// Routes //////////////////////////////////////////////////////////////////////////

//
//
//
func socketService(hub *Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}
}

//
//
//
func home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Measurements"}
	render(w, "index", context)
}

//
//
//
func staticHandler(w http.ResponseWriter, req *http.Request) {

	static_file := req.URL.Path[len(STATIC_URL):]

	if len(static_file) != 0 {

		if f, err := http.Dir(STATIC_ROOT).Open(static_file); err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}

	}
	http.NotFound(w, req)
}
