package web

import (
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/DSchalla/Claptrap/rules"
	"net/http"
	"fmt"
	"github.com/mattermost/mattermost-server/mlog"
		"github.com/DSchalla/Claptrap/analysis"
	"html/template"
	"os"
	"path/filepath"
	"path"
	"time"
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

type PageContext struct {
	URL string
	Data interface{}
}

func (s *Server) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	mlog.Debug(fmt.Sprintf("User requested resource: %s", req.URL.Path))
	s.router.HandleHTTP(w, req)
}

func (s *Server) IndexHandler(w http.ResponseWriter, req *http.Request) {
	t := s.getTemplate()
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/index.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][IndexHandler] Error parsing index template: %s", err))
	}

	foo := struct {
		Test string
	}{
		"Hello World",
	}

	s.execTemplate(t, w, foo)
}


func (s *Server) AuditHandler(w http.ResponseWriter, req *http.Request) {
	t := s.getTemplate()
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/audit.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][AuditHandler] Error parsing index template: %s", err))
	}

	events, err := s.audit.GetEvents(time.Now())

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][AuditHandler] Error getting audit events: %s", err))
	}

	data := struct {
		Events []analysis.AuditMessage
		Date string
	}{
		events,
		time.Now().Format("2006-01-02"),
	}

	mlog.Debug(fmt.Sprintf("%+v\n", data.Events))

	context := PageContext{
		URL: req.URL.Path,
		Data: data,
	}

	s.execTemplate(t, w, context)
}

func (s *Server) CaseNewHandler(w http.ResponseWriter, req *http.Request) {
	t := s.getTemplate()
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/case_new.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][AuditHandler] Error parsing index template: %s", err))
	}

	context := PageContext{
		URL: req.URL.Path,
	}

	s.execTemplate(t, w, context)
}

func (s *Server) CaseNewHandlerCreate(w http.ResponseWriter, req *http.Request) {

}

func (s *Server) execTemplate(t *template.Template, w http.ResponseWriter, context interface{}) error{
	err := t.ExecuteTemplate(w, "base", context)

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][INDEXHANDLER] Error Executing Template: %s", err))
		return err
	}

	return nil
}

func (s *Server) getTemplate() *template.Template {
	t := template.New("")
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/partials/base.html.tpl"), path.Join(s.getBasePath(), "static/partials/sidebar.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][getTemplate] Error parsing base templates: %s", err))
	}

	return t
}

func (s *Server) getBasePath() string{
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}