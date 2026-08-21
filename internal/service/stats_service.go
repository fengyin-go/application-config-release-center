package service

type ConfigItemCountStat struct {
	AppID string `json:"app_id"`
	EnvID string `json:"env_id"`
	Count int    `json:"count"`
}

func (s *Service) StatsConfigItemCountByAppEnv() []ConfigItemCountStat {
	all := s.store.ListConfigItems()
	m := make(map[string]int)
	for _, c := range all {
		key := c.AppID + "|" + c.EnvID
		m[key]++
	}
	var res []ConfigItemCountStat
	for key, cnt := range m {
		parts := split2(key, "|")
		res = append(res, ConfigItemCountStat{AppID: parts[0], EnvID: parts[1], Count: cnt})
	}
	return res
}

type ReleaseStatusStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (s *Service) StatsReleaseByStatus() []ReleaseStatusStat {
	all := s.store.ListReleases()
	m := make(map[string]int)
	for _, r := range all {
		m[r.Status]++
	}
	var res []ReleaseStatusStat
	for st, cnt := range m {
		res = append(res, ReleaseStatusStat{Status: st, Count: cnt})
	}
	return res
}

type AuditActionStat struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}

func (s *Service) StatsAuditByAction() []AuditActionStat {
	all := s.store.ListAuditLogs()
	m := make(map[string]int)
	for _, a := range all {
		m[a.Action]++
	}
	var res []AuditActionStat
	for ac, cnt := range m {
		res = append(res, AuditActionStat{Action: ac, Count: cnt})
	}
	return res
}

func split2(s, sep string) [2]string {
	var res [2]string
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			res[0] = s[:i]
			res[1] = s[i+1:]
			return res
		}
	}
	res[0] = s
	return res
}
