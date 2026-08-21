package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerReleaseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/releases", s.createRelease)
	mux.HandleFunc("GET /api/releases", s.listReleases)
	mux.HandleFunc("GET /api/releases/{id}", s.getRelease)
	mux.HandleFunc("PUT /api/releases/{id}/status", s.updateReleaseStatus)
	mux.HandleFunc("DELETE /api/releases/{id}", s.deleteRelease)
}

type createReleaseRequest struct {
	AppID   string `json:"app_id"`
	EnvID   string `json:"env_id"`
	Version string `json:"version"`
	Remark  string `json:"remark"`
	Status  string `json:"status"`
}

func (s *Server) createRelease(w http.ResponseWriter, r *http.Request) {
	var req createReleaseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rel, err := s.svc.CreateRelease(model.Release{
		AppID: req.AppID, EnvID: req.EnvID, Version: req.Version, Remark: req.Remark, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rel)
}

func (s *Server) listReleases(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReleaseFilter{
		AppID:  r.URL.Query().Get("app_id"),
		EnvID:  r.URL.Query().Get("env_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListReleases(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRelease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rel, err := s.svc.GetRelease(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rel)
}

type updateReleaseStatusRequest struct {
	Status   string `json:"status"`
	Operator string `json:"operator"`
}

func (s *Server) updateReleaseStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReleaseStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rel, err := s.svc.UpdateReleaseStatus(id, req.Status, req.Operator)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rel)
}

func (s *Server) deleteRelease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRelease(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
