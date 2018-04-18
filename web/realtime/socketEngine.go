package realtime

import (
	"fmt"
	"golang.org/x/net/websocket"
	"encoding/json"
	. "conferenceBot/config"
	"conferenceBot/database"
	"conferenceBot/game"
)

/**
 * Интерфейс Задачника
 */
type ServiceEngine struct {
	runing bool
	clients []*websocket.Conn
	botUsers map[int64]User
	db *database.Connect
	game *game.GamePlay
}

func NewServiceEngine(db *database.Connect,game *game.GamePlay) *ServiceEngine {
	engine := &ServiceEngine{runing:false,clients:[]*websocket.Conn{},botUsers:make(map[int64]User),db:db,game:game}
	return engine
}

func (self *ServiceEngine) NewService(ws *websocket.Conn)  {
	self.clients = append(self.clients,ws)
	self.Start(ws)
}
func (self *ServiceEngine) Start(ws *websocket.Conn) {
	var err error
	var receipt = ""
	var reply = Emit{}
	var resp = Emit{}
	fmt.Println("[realtime-info] Запуск Веб-сокета выполнен успешно")
	for {
		if err = websocket.Message.Receive(ws, &receipt); err != nil {
			fmt.Println("[realtime-info] Нет подключенных устройств",receipt,err)
			self.DeleteService(ws)
			break
		}
		fmt.Println("[realtime-info] Сообщение",receipt)
		err := self.readJson(receipt,&reply)
		if err != nil {
			fmt.Println("[realtime-error] Результат ошибки",receipt,reply,err)
			break
		}
		for _,client := range self.clients {
			fmt.Println("[realtime-info] Произошло событие",reply.Action)
			if self.checkUser(ws,client,reply.User) {
				if reply.Action == GET_USERS {
					resp.Action = RESPONSE_USER
					resp.Body.Message = "К серверу уже подключены \n"
					resp.User = reply.User
					for _,v := range self.botUsers {
						resp.Body.Message += v.Name + "\n"
					}
					resp.Body.Message += "Еще что нибудь?\n"
					r,_ := json.Marshal(resp)
					if err = websocket.Message.Send(client, string(r)); err != nil {
						fmt.Println("[realtime-info] Не возможно отправить клиенту",client)
						break
					}
				}
				continue
			}else{
				self.SendWeb(client,reply)
			}
		}
	}
}

func (self *ServiceEngine) SendAll(action string, body interface{}) {
	var err error
	var reply = struct{
		Action string `json:"action"`
		Body interface{} `json:"body"`
	}{}
	for _,client := range self.clients {
		reply.Action = action
		reply.Body = body
		r,_ := json.Marshal(reply)
		fmt.Println("[realtime-info] Отправлено клиенту : ",string(r))
		if err = websocket.Message.Send(client, string(r)); err != nil {
			fmt.Println("[realtime-info] Не возможно отправить клиенту",client)
			break
		}
	}
}

func (self *ServiceEngine) SendUser(user User,description string) {
	var err error
	var reply = struct{
		Action string `json:"action"`
		Body interface{} `json:"body"`
	}{}
	for _,client := range self.clients {
		reply.Action = USER_NOTIFICATION
		reply.Body = struct {
			User User `json:"user"`
			Description string `json:"description"`
		}{User: user,Description:description}
		r,_ := json.Marshal(reply)
		fmt.Println("[realtime-info] Отправлено клиенту : ",string(r))
		if err = websocket.Message.Send(client, string(r)); err != nil {
			fmt.Println("[realtime-info] Не возможно отправить клиенту",self.clients)
			break
		}
	}
}


func (self *ServiceEngine) SendWeb (ws *websocket.Conn,reply Emit) {
	if reply.Action == SEND_WEB_THREE_WINNER {
		reply.Body.Data = self.MergeUser(reply.User)
		r,_ := json.Marshal(reply)
		if err := websocket.Message.Send(ws, string(r)); err != nil {
			fmt.Println("[realtime-info] Не возможно отправить клиенту",ws)
			return
		}
	}
	if reply.Action == SEND_WEB_TWO_WINNER {
		reply.Body.Data = self.MergeUser(reply.User)
		r,_ := json.Marshal(reply)
		if err := websocket.Message.Send(ws, string(r)); err != nil {
			fmt.Println("[realtime-info] Не возможно отправить клиенту",ws)
			return
		}
	}
	if reply.Action == SEND_WEB_ONE_WINNER {
		reply.Body.Data = self.MergeUser(reply.User)
		r,_ := json.Marshal(reply)
		if err := websocket.Message.Send(ws, string(r)); err != nil {
			fmt.Println("[realtime-info] Не возможно отправить клиенту",ws)
			return
		}
	}
	users := []FullUser{}
	for _,v := range self.botUsers{
		users = append(users,self.MergeUser(v))
	}
	reply.Action = SEND_USERS
	reply.Body.Data = users
	r,_ := json.Marshal(reply)
	if err := websocket.Message.Send(ws, string(r)); err != nil {
		fmt.Println("[realtime-info] Не возможно отправить клиенту",ws)
	}
}

func (self *ServiceEngine) DeleteService(ws *websocket.Conn) {
	for i, v := range self.clients {
		if v == ws {
			copy(self.clients[i:], self.clients[i + 1:])
			self.clients[len(self.clients) - 1] = &websocket.Conn{}
			self.clients = self.clients[:len(self.clients) - 1]
		}
	}
}

func (self *ServiceEngine) readJson(body string,res interface{}) error {
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return err
	}
	return nil
}

func (self *ServiceEngine) checkUser(ws *websocket.Conn,client *websocket.Conn,user User) bool {

	if _, ok := self.botUsers[user.Chat_ID]; ok {
		if ws == client {
			return true
		}

	}
	self.botUsers[user.Chat_ID] = user
	if ws == client {
		return true
	}
	return false
}

func (self *ServiceEngine) MergeUser(u User) FullUser {
	ug,_ := self.game.FindPlayer(u.Chat_ID)
	return FullUser{
		u.Chat_ID,
		u.Phone,
		u.Name,
		u.Login,
		u.Photo,
		ug.Id,
		ug.CountErr,
		ug.CountRight,
		ug.Status,
		ug.Done,
	}
}