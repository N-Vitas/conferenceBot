package config

const (
	DATABASE_FILE 				= "./database/conference.db"
	PORT                        =":8070"
	WEB_SOCKET_ORIGIN           = "http://localhost"+PORT+"/"
	WEB_SOCKET_URL              = "ws://localhost"+PORT+"/realtime"
	TELEGRAM_BOT_API_KEY        = "559993204:AAEOyAasePH1kZ3Ba_K6wu5-lgYpqahUjnU"
	TELEGRAM_BOT_UPDATE_OFFSET  = 0
	TELEGRAM_BOT_UPDATE_TIMEOUT = 64
	GAME_LIMIT_QUESTION = 10
)