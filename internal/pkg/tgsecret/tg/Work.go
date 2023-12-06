package tg

import (
	"fmt"
	"strings"
	"time"

	JsonBase "github.com/RB-PRO/labexp/internal/pkg/tgsecret/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Отправить уведомление админам
func (TG *Telegram) Message(str string) error {
	msg := tgbotapi.NewMessage(TG.ChatNotificationID, str)
	_, err := TG.Send(msg)
	return err
}
func (TG *Telegram) Watch(db *JsonBase.Base) error {
	// Создайте новую структуру конфигурации обновления со смещением 0.
	// Смещения используются для того, чтобы убедиться, что Telegram знает,
	// что мы обработали предыдущие значения, и нам не нужно их повторять.
	updateConfig := tgbotapi.NewUpdate(0)

	// Сообщите Telegram, что мы должны ждать обновления до 30 секунд при каждом запросе.
	// Таким образом, мы можем получать информацию так же быстро,
	// как и при выполнении множества частых запросов,
	// без необходимости отправлять почти столько же.
	updateConfig.Timeout = 30

	// Начните опрос Telegram на предмет обновлений
	updates := TG.GetUpdatesChan(updateConfig)

	// Мапа, которая обеспечивает путь поиска по каждому пользователю
	userSeach := make(map[string]string)

	// map:=map[string]string
	// Давайте рассмотрим каждое обновление, которое мы получаем от Telegram
	for update := range updates {
		// Telegram может отправлять множество типов обновлений в зависимости от того,
		// чем занимается ваш бот.
		// Пока мы хотим просмотреть только сообщения,
		// чтобы исключить любые другие обновления
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			// Create a new MessageConfig. We don't have text yet,
			// so we leave it empty.

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help", "help@Geo987_bot":
				msg.Text = "/key - работа с ключами"
			case "add", "add@Geo987_bot":
				key := strings.ReplaceAll(update.Message.Text, "/add ", "")
				key = strings.TrimSpace(key)
				if key != "" && strings.Contains(key, "/") {
					db.Data.Kkeys[key] = JsonBase.Info{Access: true}
					db.Save()
					msg.Text = "Добавил и активировал ключ " + key
				} else {
					msg.Text = "Некорректный ввод, и пожалуйста не используй слэш для ключа\n/add [key]"
				}
				msg.ReplyToMessageID = update.Message.MessageID
			case "ping", "ping@Geo987_bot":
				msg.Text = time.Now().Format("2006-01-02T15:04:05")
			case "key", "key@Geo987_bot":
				msg.Text = "Выбери ключ:"
				msg.ReplyMarkup = Menu(MapKeys(db.Data))
			default:
				msg.Text = "I don't know that command"
			}
			TG.Send(msg)
			continue
		}

		if update.Message.ReplyToMessage != nil {
			// fmt.Println("ReplyToMessage", update.Message.ReplyToMessage.Text)
			switch update.Message.ReplyToMessage.Text {
			case "Выбери ключ:":
				if _, ok := db.Data.Kkeys[update.Message.Text]; ok {
					userSeach[update.Message.From.UserName] = update.Message.Text
					msg.Text = "Что мне сделать с этим ключом?"
					msg.ReplyMarkup = EditMenu()
				} else {
					msg.Text = "Я не нашёл ключ " + update.Message.Text
				}
			case "Что мне сделать с этим ключом?":
				// fmt.Println(update.Message.Text)
				switch update.Message.Text {
				case "Разрешить":
					if key, ok := userSeach[update.Message.From.UserName]; ok {
						fmt.Printf("Разрешить: Пользователь %s найден. Он искал %s\n", update.Message.From.UserName, key)
						UpdateData := db.Data.Kkeys[key]
						UpdateData.Access = true
						db.Data.Kkeys[key] = UpdateData
						msg.ReplyToMessageID = update.Message.MessageID
						msg.Text = fmt.Sprintf("Разрешил ключ %s", key)
					} else {
						msg.Text = "Разрешить: Я вообще не понял, что ты мне хочешь сказать, потому что ты не выбрал ключ\n/key"
						msg.ReplyToMessageID = update.Message.MessageID
					}
					db.Save()
				case "Заретить":
					if key, ok := userSeach[update.Message.From.UserName]; ok {
						fmt.Printf("Заретить: Пользователь %s найден. Он искал %s\n", update.Message.From.UserName, key)
						UpdateData := db.Data.Kkeys[key]
						UpdateData.Access = false
						db.Data.Kkeys[key] = UpdateData
						msg.ReplyToMessageID = update.Message.MessageID
						msg.Text = fmt.Sprintf("Запретил ключ %s", key)
					} else {
						msg.Text = "Заретить: Я вообще не понял, что ты мне хочешь сказать, потому что ты не выбрал ключ\n/key"
						msg.ReplyToMessageID = update.Message.MessageID
					}
					db.Save()
				case "Удалить":
					if key, ok := userSeach[update.Message.From.UserName]; ok {
						fmt.Printf("Удалить: Пользователь %s найден. Он искал %s\n", update.Message.From.UserName, key)
						delete(db.Data.Kkeys, key)
						msg.ReplyToMessageID = update.Message.MessageID
						msg.Text = fmt.Sprintf("Удалил ключ %s", key)
					} else {
						msg.Text = "Удалить: Я вообще не понял, что ты мне хочешь сказать, потому что ты не выбрал ключ\n/key"
						msg.ReplyToMessageID = update.Message.MessageID
					}
					db.Save()

				case "История":
					if key, ok := userSeach[update.Message.From.UserName]; ok {
						fmt.Printf("Удалить: Пользователь %s найден. Он искал %s\n", update.Message.From.UserName, key)
						msg.Text = fmt.Sprintf("Ключ %s - %v\n", key, db.Data.Kkeys[key].Access)
						for _, visit := range db.Data.Kkeys[key].VisitHistory {
							msg.Text += fmt.Sprintf(">Дата входа: %s\n", visit.Date.Format("15:04 02.01.2006"))
							msg.Text += fmt.Sprintf("--Названеи файла: %s\n", visit.FileName)
							msg.Text += fmt.Sprintf("--Пользователь ПК: %s\n", visit.UserPC)
							msg.Text += fmt.Sprintf("--IP: %s\n", visit.IP)
						}
					} else {
						msg.Text = "История: Я вообще не понял, что ты мне хочешь сказать, потому что ты не выбрал ключ\n/key"
						msg.ReplyToMessageID = update.Message.MessageID
					}
				}
			}
			TG.Send(msg)
			fmt.Printf("%+v\n", db.Data.Kkeys)
			fmt.Printf("%+v\n\n", userSeach)
			continue
		}

	}
	return nil
}

func MapKeys(datas JsonBase.Content) []string {
	mk := make([]string, 0, len(datas.Kkeys))
	for k := range datas.Kkeys {
		mk = append(mk, k)
	}
	return mk
}

// Получить значение нижнего бара для отправки сообщения
func Menu(strs []string) (key tgbotapi.ReplyKeyboardMarkup) {
	var j, indexRow int
	Cols := 2                        // к-во колонок
	Rows := CoutRow(len(strs), Cols) // к-во строк
	Col := make([][]tgbotapi.KeyboardButton, Rows)
	for i := 0; i < len(strs); i += Cols {
		j += Cols
		if j > len(strs) {
			j = len(strs)
		}
		// Row := []interface{}{}
		Row := make([]tgbotapi.KeyboardButton, Cols)
		for _, pc := range strs[i:j] {
			Row = append(Row, tgbotapi.NewKeyboardButton(pc))
		}

		Col[indexRow] = Row
		indexRow++
	}
	key.Keyboard = Col
	return key
}

func EditMenu() (key tgbotapi.ReplyKeyboardMarkup) {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Разрешить")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Заретить")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Удалить")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("История")),
	)
}

// Получить количество строк из требований количества колонок
// при заданном количестве элементов
func CoutRow(a int, Cols int) int {
	b := a / Cols
	if a-b*Cols != 0 {
		return b + 1
	}
	return b
}
