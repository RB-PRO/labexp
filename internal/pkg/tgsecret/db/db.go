package JsonBase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Структура лёгковесной БД в виде json :d
type Base struct {
	FileName string
	Data     Content
}

type Content struct {
	Keys  []string        `json:"keys"`
	Kkeys map[string]Info `json:"kkeys"`
}
type Info struct {
	Access       bool    // Доступен ли для пользователя ключ
	VisitHistory []Visit // История входов пользователя
}
type Visit struct {
	Date     time.Time // Дата попытки входа
	UserPC   string    // Имя пользователя ПК
	FileName string    // Название файла
	IP       string    // IP Адресс
}

func NewBase(FileName string) (*Base, error) {
	bs := &Base{}
	bs.FileName = FileName
	if _, err := os.Stat(FileName); !errors.Is(err, os.ErrNotExist) {
		// Прочитать файл
		fileBytes, Err := os.ReadFile(FileName)
		if Err != nil {
			return nil, fmt.Errorf("os.ReadFile: %v", Err)
		}

		// Распарсить
		Err = json.Unmarshal(fileBytes, &bs.Data)
		if Err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %v", Err)
		}

	}
	return bs, nil
}

func (db *Base) Save() error {
	jsonData, err := json.Marshal(db.Data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}
	err = os.WriteFile(db.FileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return nil
}
