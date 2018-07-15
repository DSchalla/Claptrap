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
	"github.com/gorilla/mux"
	"strings"
)

func NewServer(api plugin.API, caseManager *rules.CaseManager, audit *analysis.AuditTrail) *Server {
	s := &Server{}
	s.api = api
	s.caseManager = caseManager
	s.audit = audit
	s.router = newRouter(s)
	s.baseTemplate = s.createBaseTemplate()
	return s
}

type Server struct {
	api         plugin.API
	caseManager *rules.CaseManager
	audit       *analysis.AuditTrail
	router      *router
	baseTemplate *template.Template
}

type PageContext struct {
	URL string
	Data interface{}
}

func (s *Server) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	authorized := false
	username := ""
	userId := req.Header.Get("Mattermost-User-Id")

	if userId != "" {
		user, _ := s.api.GetUser(userId)
		username = user.Username
		authorized = strings.Contains(user.Roles, "admin")
	}

	if !authorized {
		ip := req.RemoteAddr

		if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		mlog.Warn("Unauthorized user tried to access admin panel",
			mlog.String("UserId", userId),
			mlog.String("IP", ip),
			mlog.String("UserAgent", req.UserAgent()),
		)
		http.Redirect(w, req, "/", 302)
		return
	}

	mlog.Debug("User requested resource", mlog.String("path", req.URL.Path), mlog.String("user", username))
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


	context := PageContext{
		URL: req.URL.Path,
		Data: data,
	}

	s.execTemplate(t, w, context)
}

func (s *Server) CasesHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	typeName := vars["type"]
	t := s.getTemplate()
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/cases.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][CasesHandler] Error parsing index template: %s", err))
	}

	cases, err := s.caseManager.GetForType(vars["type"])

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][CasesHandler] Error getting audit events: %s", err))
	}

	var templateCases []interface{}
	for _, engineCase := range cases {
		templateCases = append(templateCases, struct{
			Name string
			NumConditions int
			NumResponses int
			Type string
		}{
			engineCase.Name,
			len(engineCase.Conditions),
			len(engineCase.Responses),
			typeName,
		})
	}

	data := struct {
		Cases []interface{}
		Type string
	}{
		templateCases,
		strings.Title(strings.Replace(typeName, "_", " ", -1)),
	}


	context := PageContext{
		URL: req.URL.Path,
		Data: data,
	}

	s.execTemplate(t, w, context)
}

func (s *Server) CaseNewHandler(w http.ResponseWriter, req *http.Request) {
	t := s.getTemplate()
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/case_form.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][AuditHandler] Error parsing index template: %s", err))
	}

	context := PageContext{
		URL: req.URL.Path,
		Data: struct {
			CaseTypes map[string]string
			ConditionOptions map[string]string
			ResponseOptions map[string]string
			Case *rules.RawCase
			DummyStruct struct{Id, Condition string}
		}{
			s.caseManager.GetCaseTypes(),
			s.caseManager.GetConditionOptions(),
			s.caseManager.GetResponseOptions(),
			nil,
			struct{Id, Condition string}{"{INDEX}", ""},
		},
	}

	s.execTemplate(t, w, context)
}

func (s *Server) CaseNewHandlerCreate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	intercept := false

	if req.FormValue("intercept") == "Yes"{
		intercept = true
	}

	rawCase := rules.RawCase{
		Name: req.FormValue("casename"),
		Intercept: intercept,
		ConditionMatching: req.FormValue("condition_matching"),
	}

	for i:= 0; i < 10; i++ {
		prefix := fmt.Sprintf("conditions[%d]", i)
		conditionType := req.FormValue(prefix+"[type]")

		if conditionType == "" {
			break
		}

		conditionValue := ""

		if conditionType == "message_contains" || conditionType == "message_starts_with" {
			conditionValue = req.FormValue(prefix+"[condition]")
		}

		rawCond := rules.RawCondition{
			CondType: conditionType,
			Condition: conditionValue,
		}
		rawCase.Conditions = append(rawCase.Conditions, rawCond)
	}

	for i:= 0; i < 10; i++ {
		prefix := fmt.Sprintf("responses[%d]", i)
		responseType := req.FormValue(prefix+"[type]")

		if responseType == "" {
			break
		}

		responseMessage := ""

		if responseType == "message_channel" {
			responseMessage = req.FormValue(prefix+"[message]")
		}

		rawResp := rules.RawResponse{
			Action: responseType,
			Message: responseMessage,
		}
		rawCase.Responses = append(rawCase.Responses, rawResp)
	}

	newCase := rules.CreateCaseFromRawCase(rawCase)
	err := s.caseManager.Add(req.FormValue("type"), newCase)

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][CaseNewHandlerCreate] Error Adding Case: %s", err))
	}

	http.Redirect(w, req, "/plugins/com.dschalla.claptrap/cases/" + req.FormValue("type"), 302)
}

func (s *Server) CasesDeleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	typeName := vars["type"]
	caseName := vars["name"]
	s.caseManager.Delete(typeName, caseName)
	http.Redirect(w, req, "/plugins/com.dschalla.claptrap/cases/" + typeName, 302)
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
	t, err := s.baseTemplate.Clone()

	if err != nil {
		mlog.Error("[CLAPTRAP][WEB][getTemplate] Error Cloning Template", mlog.Err(err))
		return nil
	}

	return t
}

func (s *Server) createBaseTemplate() *template.Template {
	t := template.New("")
	t, err := t.ParseFiles(path.Join(s.getBasePath(), "static/partials/base.html.tpl"), path.Join(s.getBasePath(), "static/partials/sidebar.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][createBaseTemplate] Error parsing base templates: %s", err))
	}

	t, err = t.ParseGlob(path.Join(s.getBasePath(), "static/partials/conditions/*.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][createBaseTemplate] Error parsing base templates: %s", err))
	}

	t, err = t.ParseGlob(path.Join(s.getBasePath(), "static/partials/responses/*.html.tpl"))

	if err != nil {
		mlog.Error(fmt.Sprintf("[CLAPTRAP][WEB][createBaseTemplate] Error parsing base templates: %s", err))
	}

	return t
}

func (s *Server) getBasePath() string{
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}