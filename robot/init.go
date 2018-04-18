package robot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "conferenceBot/config"
	"log"
	"conferenceBot/robot/realtime"
	"fmt"
	"conferenceBot/database"
	"conferenceBot/game"
)

type Bot struct {
	API                   *tgbotapi.BotAPI        // API телеграмма
	Updates               tgbotapi.UpdatesChannel // Канал обновлений
	ActiveContactRequests []int64
	ActiveNameRequests []int64
	BotUsers map[int64]User
	Socket *realtime.BotEngine
	Game *game.GamePlay
	Db *database.Connect
	ActiveQuestionData map[int64]GameQuestion
	//BotUsers map[int64]bool
	// ID чатов, от которых мы ожидаем номер
}

func (b *Bot) Init(db *database.Connect,game *game.GamePlay) {
	b.Db = db
	b.Game = game
	b.Game.CreateGamePlay()
	b.ActiveQuestionData = make(map[int64]GameQuestion)
	botAPI, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_API_KEY)  // Инициализация API
	if err != nil {
		log.Fatal("NewBotAPI",err,TELEGRAM_BOT_API_KEY)
	}
	b.API = botAPI
	botUpdate := tgbotapi.NewUpdate(TELEGRAM_BOT_UPDATE_OFFSET)  // Инициализация канала обновлений
	botUpdate.Timeout = TELEGRAM_BOT_UPDATE_TIMEOUT
	botUpdates, err := b.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal("GetUpdatesChan",err)
	}
	b.Updates = botUpdates
	b.BotUsers = b.Db.GetUsers()
	b.Socket = realtime.NewBotEngine(botAPI)
	go b.Socket.ReadMessage()
	fmt.Println("[robot-info] Инициализация Бота выполнена успешно ")
}