package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"os"
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
	exe, _ := os.Executable()
	mlog.Debug(fmt.Sprintf("Setting Static Directory to: %s", path.Join(filepath.Dir(exe), "/static")))
	r.router.HandleFunc("/config", s.ConfigHandler).Methods("POST")
	r.router.PathPrefix("/").Handler(http.FileServer(http.Dir(path.Join(filepath.Dir(exe), "static"))))
}

func (r *router) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
