package robot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"strconv"
	"conferenceBot/config"
)

func (b *Bot) callbackQuery(callback *tgbotapi.CallbackQuery) {
	// Получаем id чата чтоб бот понимал кому слать ответ
	chatID := callback.Message.Chat.ID
	messID := callback.Message.MessageID
	answer := ""
	if question,ok := b.ActiveQuestionData[chatID];ok{
		u,ok := b.Game.FindPlayer(chatID)
		if ok == false {
			return
		}
		for _,v := range u.GetAnswersId(){
			id, ok := strconv.Atoi(v)
			if ok == nil {
				if id == int(question.Id) {
					msg := tgbotapi.NewMessage(chatID, "Вы уже отвечали на этот вопрос\n")
					msg.ParseMode = "Markdown"
					tgbotapi.NewRemoveKeyboard(false)
					b.API.Send(msg)
					return
				}
			}
		}
		answer = b.getCallbackOption(callback,question.Option1,question.Option2,question.Option3,question.Option4)
		if b.Game.CheckQuestion(question.Id,chatID,answer){
			answer = "Все верно. Ответ " + question.Answer +" "+ b.Game.LeftCounrPlayer(chatID)+"\n"
		} else {
			answer = "Сожалею ответ не верный. Правильный ответ " + question.Answer + "\n"+"*Начнем с начала.*\n" +b.Game.LeftCounrPlayer(chatID)+"\n"
		}
	}
	b.EmitWinner()
	b.createQuestion(chatID,messID,answer)
}

func (b *Bot) getCallbackOption(callback *tgbotapi.CallbackQuery,o1 string,o2 string,o3 string,o4 string) string {
	if strings.Index(callback.Data, o1) != -1 { return o1 }
	if strings.Index(callback.Data, o2) != -1 { return o2 }
	if strings.Index(callback.Data, o3) != -1 { return o3 }
	if strings.Index(callback.Data, o4) != -1 { return o4 }
	return "Fail"
}

func (b *Bot) EmitWinner() {
	if b.Game.FirstUser > 0 && b.Game.LastUser > 0 && b.Game.SecondUser > 0 {
		b.Socket.SendMessage(config.SEND_WEB_THREE_WINNER,"Третье место",b.BotUsers[b.Game.LastUser])
		return
	}
	if b.Game.FirstUser > 0 && b.Game.SecondUser > 0 {
		b.Socket.SendMessage(config.SEND_WEB_TWO_WINNER,"Второе место",b.BotUsers[b.Game.SecondUser])
		return
	}
	if b.Game.FirstUser > 0 {
		b.Socket.SendMessage(config.SEND_WEB_ONE_WINNER,"Первое место",b.BotUsers[b.Game.FirstUser])
		return
	}
}