package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerApplicationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/applications", s.createApplication)
	mux.HandleFunc("GET /api/applications", s.listApplications)
	mux.HandleFunc("GET /api/applications/{id}", s.getApplication)
	mux.HandleFunc("PUT /api/applications/{id}", s.updateApplication)
	mux.HandleFunc("DELETE /api/applications/{id}", s.deleteApplication)
	mux.HandleFunc("GET /api/applications/{id}/environments", s.getApplicationEnvironments)
	mux.HandleFunc("GET /api/applications/{id}/config-items", s.getApplicationConfigItems)
	mux.HandleFunc("GET /api/applications/{id}/releases", s.getApplicationReleases)
}

type createApplicationRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createApplication(w http.ResponseWriter, r *http.Request) {
	var req createApplicationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateApplication(model.Application{Name: req.Name, Code: req.Code, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listApplications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ApplicationFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListApplications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetApplication(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateApplicationRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) updateApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateApplicationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateApplication(id, model.Application{Name: req.Name, Code: req.Code, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteApplication(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getApplicationEnvironments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.GetApplicationEnvironments(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getApplicationConfigItems(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.GetApplicationConfigItems(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getApplicationReleases(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.GetApplicationReleases(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
