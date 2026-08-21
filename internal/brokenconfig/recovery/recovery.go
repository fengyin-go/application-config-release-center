package recovery

import (
	"configcenter/internal/brokenconfig/builder"
	"configcenter/internal/brokenconfig/model"
	"fmt"
)

type Recoverer struct{ Builder builder.Builder }

// Capture builds into an isolated copy of the caller's bundle so that a failing
// build can mutate the working copy without ever touching the stable input.
// On failure the recovered result is discarded: callers get (nil, err) and the
// passed target is left untouched.
func (r Recoverer) Capture(target *model.Bundle, fail bool) (result *model.Bundle, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = nil
			err = fmt.Errorf("build failed: %v", recovered)
		}
	}()
	result = clone(target)
	r.Builder.Build(result, fail)
	return result, nil
}

// clone returns an independent copy of b so a panicking build cannot pollute
// the caller's stable bundle through shared map state.
func clone(b *model.Bundle) *model.Bundle {
	if b == nil {
		return &model.Bundle{Values: map[string]string{}}
	}
	values := make(map[string]string, len(b.Values))
	for k, v := range b.Values {
		values[k] = v
	}
	return &model.Bundle{Values: values, Complete: b.Complete}
}
