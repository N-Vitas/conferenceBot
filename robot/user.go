package robot

import (
	. "conferenceBot/config"
)


func (b *Bot) createUser(chatID int64,login string,photo string) {
	user := User{chatID, "","",login,photo,false}
	b.Db.CreateUser(user)
	b.BotUsers[chatID] = user
	b.requestContact(chatID)
}

func (b *Bot) updateUserPhone(phone string, chatID int64) {
	for _, v := range b.BotUsers {
		if v.Chat_ID == chatID {
			v.Phone = phone
			b.Db.UpdateUser(v)
			b.BotUsers[chatID] = v
			break
		}
	}
}
func (b *Bot) updateUserName(name string, chatID int64) {
	for _, v := range b.BotUsers {
		if v.Chat_ID == chatID {
			v.Name = name
			v.Status = true
			b.Db.UpdateUser(v)
			b.BotUsers[chatID] = v
			break
		}
	}
}
