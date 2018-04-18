package main

import (
	"conferenceBot/robot"
	"github.com/bmizerany/pat"
	"net/http"
	. "conferenceBot/config"
	"conferenceBot/web"
	"conferenceBot/web/realtime"
	"golang.org/x/net/websocket"
	"fmt"
	"conferenceBot/database"
	"conferenceBot/game"
)

var conferenceBot robot.Bot

func main() {
	// Для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	db := database.NewConnect()
	game := game.NewGamePlay(db);
	web := web.Service{realtime.NewServiceEngine(db,game),db}
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/realtime", websocket.Handler(web.Realtime.NewService))
	mux := pat.New()
	mux.Get("/:page", http.HandlerFunc(web.PostHandler))
	mux.Get("/:page/", http.HandlerFunc(web.PostHandler))
	mux.Get("/", http.HandlerFunc(web.PostHandler))
	http.Handle("/", mux)

	fmt.Println("[main-info] Сервис конференц бота. Работает на порту "+PORT)
	go http.ListenAndServe(PORT, nil)
	// После запуска сокета в потоке запускаем бота
	conferenceBot.Init(db,game)
	conferenceBot.Start() // Бот остается в основном потоке
}
