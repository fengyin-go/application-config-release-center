package handler_test

import (
	"configcenter/internal/brokenconfig/builder"
	"configcenter/internal/brokenconfig/cache"
	"configcenter/internal/brokenconfig/handler"
	"configcenter/internal/brokenconfig/model"
	"configcenter/internal/brokenconfig/recovery"
	"errors"
	"reflect"
	"testing"
)

func TestRecoveredBuildNeverPublishesPartialBundle(t *testing.T) {
	recoverer := recovery.Recoverer{Builder: builder.Builder{}}
	target := &model.Bundle{Values: map[string]string{"stable": "yes"}, Complete: true}
	result, err := recoverer.Capture(target, true)
	if err == nil || result != nil {
		t.Errorf("recovery returned result=%v err=%v", result, err)
	}
	if !reflect.DeepEqual(target.Values, map[string]string{"stable": "yes"}) {
		t.Errorf("existing bundle was changed to %v", target.Values)
	}
	c := cache.New()
	c.Put("partial", &model.Bundle{Values: map[string]string{"partial": "visible"}})
	if _, ok := c.Get("partial"); ok {
		t.Error("cache accepted an incomplete bundle")
	}
	h := handler.Handler{Recoverer: recoverer, Cache: cache.New()}
	if err := h.Load("broken", &model.Bundle{Values: map[string]string{}}, true); err == nil {
		t.Error("handler hid build failure")
	}
	if value, ok := h.Read("broken"); ok {
		t.Errorf("later read exposed partial bundle %+v", value)
	}
	complete := &model.Bundle{Values: map[string]string{"ready": "yes"}, Complete: true}
	if err := h.Publish("errored", complete, errors.New("source failed")); err == nil {
		t.Error("handler hid source failure")
	}
	if value, ok := h.Read("errored"); ok {
		t.Errorf("handler cached an errored result %+v", value)
	}
}
