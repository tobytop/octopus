package web

import (
	"net/http"
	"octopus/metadata"
	"octopus/service"
)

type WebService struct {
	port          string
	routerservice *service.RouterService
}

func NewWebService(port string, routerservice *service.RouterService) *WebService {
	return &WebService{
		routerservice: routerservice,
		port:          port}
}

func (web *WebService) StartWebService() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/protoup", web.upprotoservice)
	http.ListenAndServe(web.port, mux)
}

func (web *WebService) upprotoservice(w http.ResponseWriter, r *http.Request) {
	//web.routerservice.BuildRegTable("example1.so")
}

func WithWebService() metadata.OptionBuilder[service.RouterService] {
	return func(rs *service.RouterService) {
		NewWebService("9009", rs).StartWebService()
	}
}
