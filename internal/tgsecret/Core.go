package secret

import (
	"fmt"

	SecretAPI "github.com/RB-PRO/labexp/internal/pkg/tgsecret/api"
	JsonBase "github.com/RB-PRO/labexp/internal/pkg/tgsecret/db"
	"github.com/RB-PRO/labexp/internal/pkg/tgsecret/tg"
)

// Структура для работы приложения для защиты файла
type Secret struct {
	api         *SecretAPI.API
	tg          *tg.Telegram
	bs          *JsonBase.Base
	messageChan chan string
}

// Создаём приложение сервиса
func New() (*Secret, error) {
	var Err error
	s := &Secret{}
	// bs := &JsonBase.Base{}
	bs, Err := JsonBase.NewBase("keys.json")
	fmt.Println("bs before assignment:", bs) // Добавьте эту строку
	if Err != nil {
		return nil, fmt.Errorf("JsonBase.NewBase: %v", Err)
	}
	if bs == nil {
		return nil, fmt.Errorf("JsonBase.NewBase returned a nil pointer")
	}
	s.bs = bs
	fmt.Println("s.bs after assignment:", s.bs) // Добавьте эту строку
	fmt.Println("Загружена БД")

	// Канал общения API и TG
	s.messageChan = make(chan string, 1)

	// API
	s.api = SecretAPI.NewAPI(s.bs, s.messageChan)

	// TELEGRAM
	cf, Err := tg.LoadConfig("tg.json")
	if Err != nil {
		return nil, fmt.Errorf("tg.LoadConfig: %v", Err)
	}
	fmt.Println("Загрузил конфиг для телеграма")

	s.tg = &tg.Telegram{}
	s.tg, Err = tg.NewTelegram(cf)
	if Err != nil {
		return nil, fmt.Errorf("tg.NewTelegram: %v", Err)
	}
	fmt.Println("Создано приложение телеграм")

	return s, nil
}

func (s *Secret) Run() error {
	go s.api.Watch()
	go s.tg.Watch(s.bs)
	for {
		select {
		case Message := <-s.messageChan:
			// fmt.Println(Message)
			s.tg.Message(Message)
		}
	}
	// return nil
}
