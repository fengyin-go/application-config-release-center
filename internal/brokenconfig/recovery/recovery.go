package recovery

import (
	"configcenter/internal/brokenconfig/builder"
	"configcenter/internal/brokenconfig/model"
	"fmt"
)

type Recoverer struct{ Builder builder.Builder }

func (r Recoverer) Capture(target *model.Bundle, fail bool) (result *model.Bundle, err error) {
	result = target
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("build failed: %v", recovered)
		}
	}()
	r.Builder.Build(target, fail)
	return result, nil
}
