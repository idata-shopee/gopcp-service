package gopcp_service

// pcp service:
//  (1) dependencies (like other services), called as Resource
//  (2) connection (writer, reader)
//  (3) sandbox provided functions

import (
	"log"
	"net/http"
	"strconv"
)

type Attachment struct {
	Conn     interface{} // connection
	Resource interface{} // global resource
	Session  interface{} // session information
}

type HttpConn struct {
	W http.ResponseWriter
	R *http.Request
}

func StartHttpServer(port int, routes []Route) error {
	router := GetRouter(routes)
	log.Println("try to start server at " + strconv.Itoa(port))
	return http.ListenAndServe(":"+strconv.Itoa(port), router)
}
