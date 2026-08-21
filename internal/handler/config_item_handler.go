package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerConfigItemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/config-items", s.createConfigItem)
	mux.HandleFunc("GET /api/config-items", s.listConfigItems)
	mux.HandleFunc("GET /api/config-items/{id}", s.getConfigItem)
	mux.HandleFunc("PUT /api/config-items/{id}", s.updateConfigItem)
	mux.HandleFunc("DELETE /api/config-items/{id}", s.deleteConfigItem)
}

type createConfigItemRequest struct {
	AppID       string `json:"app_id"`
	EnvID       string `json:"env_id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
	Description string `json:"description"`
	Encrypted   bool   `json:"encrypted"`
	Status      string `json:"status"`
}

func (s *Server) createConfigItem(w http.ResponseWriter, r *http.Request) {
	var req createConfigItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateConfigItem(model.ConfigItem{
		AppID: req.AppID, EnvID: req.EnvID, Key: req.Key, Value: req.Value,
		ValueType: req.ValueType, Description: req.Description, Encrypted: req.Encrypted, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listConfigItems(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConfigItemFilter{
		AppID:   r.URL.Query().Get("app_id"),
		EnvID:   r.URL.Query().Get("env_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListConfigItems(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getConfigItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetConfigItem(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateConfigItemRequest struct {
	AppID       string `json:"app_id"`
	EnvID       string `json:"env_id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
	Description string `json:"description"`
	Encrypted   bool   `json:"encrypted"`
	Status      string `json:"status"`
}

func (s *Server) updateConfigItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateConfigItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateConfigItem(id, model.ConfigItem{
		AppID: req.AppID, EnvID: req.EnvID, Key: req.Key, Value: req.Value,
		ValueType: req.ValueType, Description: req.Description, Encrypted: req.Encrypted, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteConfigItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteConfigItem(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
