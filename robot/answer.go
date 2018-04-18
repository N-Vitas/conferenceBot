package robot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "conferenceBot/config"
)

// Запросить номер телефона
func (b *Bot) requestContact(chatID int64) {
	// Готовим сообщение
	rcm := tgbotapi.NewMessage(chatID, "Приветствую. Согласны ли вы предоставить ваш номер телефона для регистрации в системе?")
	// Кнопка для положительного ответа
	yes := tgbotapi.NewKeyboardButtonContact("Да")
	// Кнопка для отрицательного ответа
	no := tgbotapi.NewKeyboardButton("Нет")
	// Собираем кнопки в ряд и добовляем в сообщение
	rcm.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{yes, no})
	// отправляем пользователю
	b.API.Send(rcm)
	// Ставим влаг что мы ждем от этого пользователя номер
	b.addContactRequestID(chatID)
}

// Проверка принятого контакта
func (b *Bot) checkRequestContactReply(update tgbotapi.Update) {
	if update.Message.Contact != nil {  // Проверяем, содержит ли сообщение контакт
		if update.Message.Contact.UserID == update.Message.From.ID {  // Проверяем действительно ли это контакт отправителя
			b.updateUserPhone(update.Message.Contact.PhoneNumber[len(update.Message.Contact.PhoneNumber)-10:],update.Message.Chat.ID)
			b.deleteContactRequestID(update.Message.Chat.ID)  // Удаляем ChatID из списка ожидания телефона
			b.addActiveNameRequests(update.Message.Chat.ID) // Ставим влаг что мы ждем от этого пользователя имя
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Как мне к вам обращаться? Напишите свое имя!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)  // Убираем клавиатуру
			b.API.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер телефона, который вы предоставили, принадлежит не вам!")
			b.API.Send(msg)
			b.requestContact(update.Message.Chat.ID)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Если вы не предоставите ваш номер телефона, вы не сможете пользоваться системой!")
		b.API.Send(msg)
		b.requestContact(update.Message.Chat.ID)
	}
}

// Знакомство с пользователем
func (b *Bot) checkRequestNameReply(update tgbotapi.Update) {
	b.updateUserName(update.Message.Text,update.Message.Chat.ID)
	b.deleteActiveNameRequests(update.Message.Chat.ID)  // Удаляем ChatID из списка ожидания имени
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приятно познакомиться "+ update.Message.Text +". Можешь звать меня просто Бот!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)  // Убираем клавиатуру
	b.API.Send(msg)
	b.Socket.SendMessage(GET_USERS,"Дай мне список пользователей",b.BotUsers[update.Message.Chat.ID])
}