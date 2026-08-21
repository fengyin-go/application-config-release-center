package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerEnvironmentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/environments", s.createEnvironment)
	mux.HandleFunc("GET /api/environments", s.listEnvironments)
	mux.HandleFunc("GET /api/environments/{id}", s.getEnvironment)
	mux.HandleFunc("PUT /api/environments/{id}", s.updateEnvironment)
	mux.HandleFunc("DELETE /api/environments/{id}", s.deleteEnvironment)
	mux.HandleFunc("GET /api/environments/{id}/config-items", s.getEnvironmentConfigItems)
	mux.HandleFunc("GET /api/environments/{id}/releases", s.getEnvironmentReleases)
}

type createEnvironmentRequest struct {
	AppID       string `json:"app_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (s *Server) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var req createEnvironmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEnvironment(model.Environment{AppID: req.AppID, Name: req.Name, Code: req.Code, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEnvironments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EnvironmentFilter{
		AppID:   r.URL.Query().Get("app_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListEnvironments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEnvironment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEnvironment(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateEnvironmentRequest struct {
	AppID       string `json:"app_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (s *Server) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEnvironmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEnvironment(id, model.Environment{AppID: req.AppID, Name: req.Name, Code: req.Code, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEnvironment(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getEnvironmentConfigItems(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.GetEnvironmentConfigItems(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getEnvironmentReleases(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	items, err := s.svc.GetEnvironmentReleases(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
