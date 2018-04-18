package robot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	. "conferenceBot/config"
	//"fmt"
)

func (b *Bot) chatUser(update tgbotapi.Update) {
	// Получаем id чата чтоб бот понимал кому слать ответ
	chatID := update.Message.Chat.ID
	// Получаем пользователя
	user := b.BotUsers[chatID]
	u,_ := b.Game.FindPlayer(chatID)
	//fmt.Println(b.BotUsers)
	// Проверяем ли пользователь номер или нет
	if len(user.Phone) > 0 && user.Status{
		switch true {
		case strings.Index(strings.ToLower(update.Message.Text), "ты дурак") != -1:
			msg := tgbotapi.NewMessage(chatID, "Сам дурак")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			b.API.Send(msg)
			return
		case strings.ToLower(update.Message.Text) == "да":
			b.Game.AddPlayer(chatID)
			b.createQuestion(chatID,0,"")
			return
		default:
			text := ""
			if u.Status == 0 {
				text = "Хочешь поиграть?"
			}else{
				text = "Продолжим игру?"
			}
			msg := tgbotapi.NewMessage(chatID, text)
			tgbotapi.NewRemoveKeyboard(false)
			b.API.Send(msg)
			b.Socket.SendMessage(SEND_USERS,"Дай мне список пользователей",b.BotUsers[update.Message.Chat.ID])
			return
		}
	} else {
		// Проверяем представился ли он или нет
		if b.findActiveNameRequests(chatID) {
			b.checkRequestNameReply(update)  // Если да -> проверяем
			return
		}
		// Если номера нет, то проверяем ждём ли мы контакт от этого ChatID
		if b.findContactRequestID(chatID) {
			b.checkRequestContactReply(update)  // Если да -> проверяем
			return
		} else {
			b.requestContact(chatID)  // Если нет -> запрашиваем его
			return
		}
	}
}

