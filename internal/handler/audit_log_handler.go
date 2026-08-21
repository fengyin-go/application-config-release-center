package handler

import (
	"net/http"

	"configcenter/internal/model"
	"configcenter/pkg/httpx"
)

func (s *Server) registerAuditLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/audit-logs", s.listAuditLogs)
}

func (s *Server) listAuditLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AuditLogFilter{
		Operator: r.URL.Query().Get("operator"),
		Action:   r.URL.Query().Get("action"),
		Resource: r.URL.Query().Get("resource"),
	}
	items, total, err := s.svc.ListAuditLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
