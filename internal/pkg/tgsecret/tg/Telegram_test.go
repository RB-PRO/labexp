package tg

import (
	"fmt"
	"testing"

	JsonBase "github.com/RB-PRO/labexp/internal/pkg/tgsecret/db"
)

// go test -v -run ^TestTG$ github.com/RB-PRO/labexp/internal/pkg/tgsecret/tg
func TestTG(t *testing.T) {

	cf, Errcf := LoadConfig("..//..//..//..//tg.json")
	if Errcf != nil {
		t.Error(Errcf)
	}
	fmt.Println("Load config telegram")

	tg, Errtg := NewTelegram(cf)
	if Errtg != nil {
		t.Error(Errtg)
	}
	fmt.Println("Create telegram")

	bs, ErrBS := JsonBase.NewBase("..//..//..//..//keys.json")
	if ErrBS != nil {
		t.Error(ErrBS)
	}
	fmt.Println("Load base")

	tg.Watch(bs)
}
