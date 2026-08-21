package builder

import "configcenter/internal/brokenconfig/model"

type Builder struct{}

func (Builder) Build(target *model.Bundle, fail bool) {
	target.Values["loaded"] = "true"
	if fail {
		target.Values["partial"] = "visible"
		panic("invalid encrypted value")
	}
	target.Complete = true
}
