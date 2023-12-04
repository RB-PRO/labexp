package api

import (
	"fmt"
	"time"

	JsonBase "github.com/RB-PRO/labexp/internal/pkg/tgsecret/db"
	"github.com/gin-gonic/gin"
)

// Структура API
type API struct {
	*gin.Engine
}

// Создать настройку для endpoint
func NewAPI(db *JsonBase.Base, c chan string) *API {
	r := setupRouter(db, c)
	return &API{r}
}

// Запускаем API
func (api *API) Watch() error {
	return api.Run(":8080")
}

func setupRouter(db *JsonBase.Base, ch chan string) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/check", func(c *gin.Context) {
		key := c.Query("key")
		ip := c.Query("ip")
		user := c.Query("user")
		file := c.Query("file")
		Message := fmt.Sprintf("Попытка входа с ip %s\nФайл: %s\nПользователь: %s", ip, user, file)
		if _, ok := db.Data.Kkeys[key]; ok {
			val := db.Data.Kkeys[key]
			val.VisitHistory = append(val.VisitHistory, JsonBase.Visit{
				Date:     time.Now(),
				IP:       ip,
				UserPC:   user,
				FileName: file,
			})
		} else {
			Message += "\nВНИМАНИЕ! Такого ключа в базе вообще нет"
		}
		ch <- Message

		if val, ok := db.Data.Kkeys[key]; ok {
			if val.Access {
				c.String(200, "1") // Работаем
			} else {
				c.String(200, "0")
			}
		} else {
			c.String(200, "0")
		}

		// close(ch)
	})
	return r
}

// func Serve(db *JsonBase.Base) {
// 	r := setupRouter(db)
// 	r.Run(":8080")
// }
