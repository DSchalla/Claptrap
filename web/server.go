package web

import (
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/DSchalla/Claptrap/rules"
	"net/http"
	"fmt"
	"github.com/mattermost/mattermost-server/mlog"
		"github.com/DSchalla/Claptrap/analysis"
)

func NewServer(api plugin.API, caseManager *rules.CaseManager, audit *analysis.AuditTrail) *Server {
	s := &Server{}
	s.api = api
	s.caseManager = caseManager
	s.audit = audit
	s.router = newRouter(s)
	return s
}

type Server struct {
	api         plugin.API
	caseManager *rules.CaseManager
	audit       *analysis.AuditTrail
	router      *router
}

func (s *Server) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	mlog.Debug(fmt.Sprintf("User requested resource: %s", req.URL.Path))
	s.router.HandleHTTP(w, req)
}

func (s *Server) ConfigHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Config Stored")
}
