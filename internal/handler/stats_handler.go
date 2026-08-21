package handler

import (
	"net/http"

	"configcenter/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/config-item-counts", s.statsConfigItemCounts)
	mux.HandleFunc("GET /api/stats/release-status", s.statsReleaseStatus)
	mux.HandleFunc("GET /api/stats/audit-actions", s.statsAuditActions)
}

func (s *Server) statsConfigItemCounts(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsConfigItemCountByAppEnv())
}

func (s *Server) statsReleaseStatus(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsReleaseByStatus())
}

func (s *Server) statsAuditActions(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsAuditByAction())
}
