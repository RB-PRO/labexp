package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	JsonBase "github.com/RB-PRO/labexp/internal/pkg/tgsecret/db"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	bs, ErrBS := JsonBase.NewBase("..//..//..//..//keys.json")
	if ErrBS != nil {
		t.Error(ErrBS)
	}
	fmt.Println("Load base")
	c := make(chan string, 1)
	router := setupRouter(bs, c)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

	reqCheck, _ := http.NewRequest("GET", "/check?key=123&ip=1.1.1.1&user=roma&file=file", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, reqCheck)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "1", w2.Body.String())

	fmt.Println(<-c)
}
