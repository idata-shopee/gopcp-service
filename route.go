package gopcp_service

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type HTTPHandler func(http.ResponseWriter, *http.Request)

type Route struct {
	URLPattern       *regexp.Regexp
	Handler          HTTPHandler
	MethodsSupported map[string]bool
}

type Router struct {
	Routes []Route
}

func NewRoute(
	pattern string,
	handler HTTPHandler,
	methods map[string]bool,
) Route {
	route := Route{}

	compiledPattern := regexp.MustCompile(pattern)
	route.URLPattern = compiledPattern
	route.Handler = handler
	route.MethodsSupported = methods

	return route
}

func (r Route) Match(req *http.Request) bool {
	match := r.URLPattern.FindStringSubmatch(req.URL.Path)

	if len(match) > 0 {
		for i, name := range r.URLPattern.SubexpNames() {
			if i != 0 && name != "" {
				req.Form.Add(name, match[i])
			}
		}
		return true
	} else {
		return false
	}
}

func (r Route) Handle(w http.ResponseWriter, req *http.Request) {
	if val, ok := r.MethodsSupported[req.Method]; ok && val == true {
		r.Handler(w, req)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unsupported method: " + req.Method))
	}
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("access %s", req.URL)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in %s", r)
		}
	}()
	for _, route := range r.Routes {
		if route.Match(req) {
			route.Handle(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func GetRouter(routes []Route) Router {
	return Router{routes}
}
