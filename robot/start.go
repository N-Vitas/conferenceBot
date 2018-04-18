package robot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"conferenceBot/config"
)

// Основной цикл бота
func (b *Bot) Start() {
	for update := range b.Updates {
		if update.Message != nil {
			// Если сообщение есть и его длина больше 0 -> начинаем обработку
			go b.analyzeUpdate(update)
		}
		if update.CallbackQuery != nil {
			// Если есть событие нажатия на ответ
			b.callbackQuery(update.CallbackQuery)
		}
	}
}

// Начало обработки сообщения
func (b *Bot) analyzeUpdate(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if b.findUser(chatID) {
		b.chatUser(update)
	} else {
		userName := update.Message.Chat.UserName
		userPhoto,_ := b.API.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(int(chatID)))
		avatar := "/static/noavatar.png"
		if userPhoto.TotalCount > 0 {
			photo := userPhoto.Photos[0][0]
			avatar,_ = b.API.GetFileDirectURL(photo.FileID)
		}
		b.createUser(chatID,userName,avatar)
	}
}

func (b *Bot) findUser(chatID int64) bool {
	if _, ok := b.BotUsers[chatID]; ok {
		return true
	}
	return false
}


func (b *Bot) addContactRequestID(chatID int64) {
	b.ActiveContactRequests = append(b.ActiveContactRequests, chatID)
}
func (b *Bot) findContactRequestID(chatID int64) bool {
	for _, v := range b.ActiveContactRequests {
		if v == chatID {
			return true
		}
	}
	return false
}
func (b *Bot) deleteContactRequestID(chatID int64) {
	for i, v := range b.ActiveContactRequests {
		if v == chatID {
			copy(b.ActiveContactRequests[i:], b.ActiveContactRequests[i + 1:])
			b.ActiveContactRequests[len(b.ActiveContactRequests) - 1] = 0
			b.ActiveContactRequests = b.ActiveContactRequests[:len(b.ActiveContactRequests) - 1]
		}
	}
}

func (b *Bot) addActiveNameRequests(chatID int64) {
	b.ActiveNameRequests = append(b.ActiveNameRequests, chatID)
}
func (b *Bot) findActiveNameRequests(chatID int64) bool {
	for _, v := range b.ActiveNameRequests {
		if v == chatID {
			return true
		}
	}
	return false
}
func (b *Bot) deleteActiveNameRequests(chatID int64) {
	for i, v := range b.ActiveNameRequests {
		if v == chatID {
			copy(b.ActiveNameRequests[i:], b.ActiveNameRequests[i + 1:])
			b.ActiveNameRequests[len(b.ActiveNameRequests) - 1] = 0
			b.ActiveNameRequests = b.ActiveNameRequests[:len(b.ActiveNameRequests) - 1]
		}
	}
}

func (b *Bot) createQuestion(chatID int64,messID int,context string) {
	b.Socket.SendMessage(config.EMIT_WEB,"Cписок пользователей",b.BotUsers[chatID])
	if b.Game.GetStatus() {
		u,_ := b.Game.FindPlayer(chatID)
		q := b.Game.GetRandomQuestion(u)
		b.ActiveQuestionData[chatID] = q
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text: q.Option1,
					CallbackData:&q.Option1,
				},
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text: q.Option2,
					CallbackData:&q.Option2,
				},
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text: q.Option3,
					CallbackData:&q.Option3,
				},
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text: q.Option4,
					CallbackData:&q.Option4,
				},
			),
		)
		if u.Status == 1 {
			if messID == 0 {
				msg := tgbotapi.NewMessage(chatID, "Вопрос такой : "+q.Question)
				msg.ParseMode = "Markdown"
				msg.ReplyMarkup = markup
				tgbotapi.NewRemoveKeyboard(false)
				b.API.Send(msg)
			} else {
				msgT := tgbotapi.NewEditMessageText(chatID,messID, context + "Следующий вопрос: "+q.Question)
				msgT.ParseMode = "Markdown"
				b.API.Send(msgT)
				msg := tgbotapi.NewEditMessageReplyMarkup(chatID,messID,markup)
				b.API.Send(msg)
			}
		}else{
			if messID == 0 {
				msg := tgbotapi.NewMessage(chatID,"Поздравляю вы победили")
				msg.ParseMode = "Markdown"
				tgbotapi.NewRemoveKeyboard(false)
				b.API.Send(msg)
			} else {
				msg := tgbotapi.NewEditMessageText(chatID, messID,"Поздравляю вы победили")
				msg.ParseMode = "Markdown"
				tgbotapi.NewRemoveKeyboard(false)
				b.API.Send(msg)
			}
		}
		return
	}
	msg := tgbotapi.NewMessage(chatID,"К сожалению игра закончена")
	msg.ParseMode = "Markdown"
	tgbotapi.NewRemoveKeyboard(false)
	b.API.Send(msg)
}

