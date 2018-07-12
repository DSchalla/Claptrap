package web

import (
	"github.com/gorilla/mux"
	"net/http"
			"path"
	"github.com/mattermost/mattermost-server/mlog"
	"fmt"
)

func newRouter(s *Server) *router {
	r := &router{}
	r.router = mux.NewRouter()
	r.setRoutes(s)
	return r
}

type router struct {
	router *mux.Router
}

func (r *router) setRoutes(s *Server) {
	base := path.Join(s.getBasePath(), "/static")
	mlog.Debug(fmt.Sprintf("Setting Static Directory to: %s", base))
	r.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(base))))
	r.router.HandleFunc("/", s.IndexHandler)
	r.router.HandleFunc("/audit", s.AuditHandler)
	r.router.HandleFunc("/case_new", s.CaseNewHandler).Methods("GET")
	r.router.HandleFunc("/case_new", s.CaseNewHandlerCreate).Methods("POST")
	r.router.HandleFunc("/cases/{type}", s.CasesHandler).Methods("GET")
	r.router.HandleFunc("/cases/{type}/{name}/delete", s.CasesDeleteHandler).Methods("POST")
}

func (r *router) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
