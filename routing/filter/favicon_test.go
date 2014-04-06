package filter

import (
	"github.com/Centny/Cny4go/routing/httptest"
	"testing"
)

func TestFavicon(t *testing.T) {
	ico := NewFavicon("favicon.ico")
	if ico == nil {
		t.Error("not new")
		return
	}
	httptest.Th(ico, "")
	httptest.Tnh(ico, "")
	ico2 := NewFavicon("favicon.ico2")
	if ico2 != nil {
		t.Error("not right")
	}
}
