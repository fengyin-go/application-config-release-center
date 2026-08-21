package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerConfigVersionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/config-versions", s.listConfigVersions)
}

func (s *Server) listConfigVersions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConfigVersionFilter{
		ConfigItemID: r.URL.Query().Get("config_item_id"),
	}
	items, total, err := s.svc.ListConfigVersions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
