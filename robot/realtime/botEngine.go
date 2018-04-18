package realtime

import (
	"encoding/json"
	"log"
	. "conferenceBot/config"
	"golang.org/x/net/websocket"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)


type BotEngine struct {
	bot *tgbotapi.BotAPI
	ws *websocket.Conn
	origin string
	url string
}

func (self *BotEngine) FillStruct(m map[string]interface{}) error {
	data, _ := json.Marshal(m)
	err := json.Unmarshal(data, &self)
	return err
}

func NewBotEngine(bot *tgbotapi.BotAPI) *BotEngine {
	engine := &BotEngine{bot:bot,origin:WEB_SOCKET_ORIGIN,url:WEB_SOCKET_URL}
	ws, err := websocket.Dial(engine.url, "", engine.origin)
	if err != nil {
		log.Fatal(err)
	}
	engine.ws = ws
	return engine
}

func (self *BotEngine) SendMessage(action string,message string,user User){
	s := Emit{Action:action,Body:Content{Message:message},User:user}
	r,_ := json.Marshal(s)
	fmt.Println("[robot-info] Подготовка к отправке ",string(r))
	if _, err := self.ws.Write(r); err != nil {
		log.Fatal(err)
	}
}
func (self *BotEngine) ReadMessage(){
	msg := make([]byte, 1024)
	message := Emit{}
	fmt.Println("[robot-info] Подключение к Веб-сокету выполнено успешно ")
	for{
		if n, err := self.ws.Read(msg); err == nil {
			body := msg[:n]
			if err := self.readJson(string(body),&message); err == nil {
				fmt.Printf("[robot-info] Событие от сервера : %s\n", message.Action)
				switch message.Action {
				case RESPONSE_USER:
					fmt.Printf("[robot-info] Ответ пользователю %s : %s\n", message.User.Chat_ID, message.Body.Message)
					msg := tgbotapi.NewMessage(message.User.Chat_ID, message.Body.Message+"!")
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)  // Убираем клавиатуру
					self.bot.Send(msg)
				default:
					fmt.Printf("[robot-info] Команда от сервера : %v\n",message.Action)
				}
			}else{
				fmt.Printf("[robot-info] Ошибка конвертации json : %v\n",err)
			}
		}else{
			fmt.Printf("[robot-info] Ошибка чтения сокета : %s\n", err)
			break
		}
	}
}

func (self *BotEngine) readJson(body string,res interface{}) error {
	err := json.Unmarshal([]byte(body),&res)
	if err != nil {
		return err
	}
	return nil
}
