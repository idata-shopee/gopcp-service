package gopcp_service

// pcp service:
//  (1) dependencies (like other services), called as Resource
//  (2) connection (writer, reader)
//  (3) sandbox provided functions

import (
	rpc "github.com/idata-shopee/gopcp_rpc"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	log.Println("try to start http server at " + strconv.Itoa(port))
	return http.ListenAndServe(":"+strconv.Itoa(port), router)
}

// blocking service
func StartTcpServer(port int, generateSandbox rpc.GenerateSandbox) error {
	log.Println("try to start tcp server at " + strconv.Itoa(port))
	if server, err := rpc.GetPCPRPCServer(port, generateSandbox); err != nil {
		return err
	} else {
		defer server.Close()

		// blocking forever
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()

		return nil
	}
}
