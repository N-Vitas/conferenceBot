package game

import (
	. "conferenceBot/config"
	"math/rand"
	"conferenceBot/database"
	"strconv"
)

type AbstractGamePlay interface {
	CreateGamePlay()
	GetLimit() int64
	SetLimit(limit int64)
	CheckLimit(countRight int64) bool
	GetRandomQuestion(gameUser *GamePlayUser) GameQuestion
	CheckQuestion(Id int64, chatId int64,answer string) bool
	AddPlayer(chatId int64)
	FindPlayer(chatId int64) GamePlayUser
	GetStatus() bool
	LeftCounrPlayer(chatId int64) string
}
type GamePlay struct {
	limit int64
	answers []GameQuestion
	users map[int64]*GamePlayUser
	db *database.Connect
	FirstUser int64
	SecondUser int64
	LastUser int64
	GameStatus bool
}

func NewGamePlay(db *database.Connect) *GamePlay {
	games := &GamePlay{
		answers:[]GameQuestion{},
		users:make(map[int64]*GamePlayUser),
		db:db,
		limit:GAME_LIMIT_QUESTION,
		FirstUser:0,
		SecondUser:0,
		LastUser:0,
		GameStatus:true,
	}
	return games
}

func (g *GamePlay) CreateGamePlay() {
	g.answers = g.db.GetGameQuestions()
}

func (g *GamePlay) GetLimit() int64{
	return g.limit
}
func (g *GamePlay) SetLimit(limit int64){

}
func (g *GamePlay) CheckLimit(countRight int64) bool  {
	return countRight >= g.limit
}
func (g *GamePlay) GetRandomQuestion(gameUser *GamePlayUser) GameQuestion  {
	answer := g.getRandomQuestion()
	done := gameUser.GetAnswersId()
	if len(done) > 0 {
		for g.checkRandom(answer,done){
			answer = g.getRandomQuestion()
		}
	}
	return answer
}
func (g *GamePlay) checkRandom(answer GameQuestion, done []string) bool {
	for _,i := range done{
		c ,_ := strconv.ParseInt(i,10,64)
		if answer.Id == c {
			return true
		}
	}
	return false
}
func (g *GamePlay) getRandomQuestion() GameQuestion {
	return g.answers[rand.Intn(len(g.answers))]
}
// Вот тут я поломал голову проверка ответа влияет на всю игру
func (g *GamePlay) CheckQuestion(Id int64, chatID int64,answer string) bool {
	// Перебераем все вопросы и иищем нужный
	if user, ok := g.users[chatID]; ok {
		for _,a := range g.answers{
			// Если вопрос найден обновляем пользователя
			if a.Id == Id {
				if a.Answer == answer {
					user.UpCountRight() // Добавляем колличество правельных ответов
					user.AddAnswerId(Id) // Добавляем вопрос в исключение чтобы он снова не появился
					if user.CountRight >= g.limit { // Если игрок ответил на все вопросы верно
						user.Status = 0
						if g.FirstUser == 0 {
							g.FirstUser = chatID
						}
						if g.FirstUser > 0 && g.FirstUser != chatID && g.SecondUser == 0{
							g.SecondUser = chatID
						}
						if g.FirstUser > 0 && g.SecondUser > 0 && g.SecondUser != chatID && g.LastUser == 0 {
							g.LastUser = chatID
						}
						if g.FirstUser > 0 && g.SecondUser > 0 && g.LastUser > 0 {
							g.GameStatus = false
						}
					}
					g.db.UpdateGamePlayer(*user)
					return true
				}
				user.UpCountErr()
				user.ClearAnswerId()
				g.db.UpdateGamePlayer(*user)
				return false
			}
		}
	}
	return false
}
func (g *GamePlay) AddPlayer(chatId int64) {
	if _,ok := g.users[chatId]; ok {
		return
	}
	user := &GamePlayUser{
		Id:0,
		Chat_id:chatId,
		CountErr:0,
		CountRight:0,
		Status:1,
	}
	user.ClearAnswerId()
	g.users[chatId] = user
	g.db.CreateGamePlayer(*user)
}
func (g *GamePlay) FindPlayer(chatId int64) (*GamePlayUser,bool) {
	if u, ok := g.users[chatId]; ok{
		return u,true
	}
	u := &GamePlayUser{
		Id:0,
		Chat_id:chatId,
		CountErr:0,
		CountRight:0,
		Status:0,
	}
	g.db.CreateGamePlayer(*u)
	return u,false
}
func (g *GamePlay) GetStatus() bool {
	return g.GameStatus
}
func (g *GamePlay) LeftCounrPlayer(chatId int64) string {
	if u,ok := g.FindPlayer(chatId); ok {
		c := g.GetLimit() - u.CountRight
		if c > 0 {
			return "Тебе осталось ответить на "+strconv.Itoa(int(c))+" вопросов"
		}
		if c <= 0 {
			return "Ты становишься победителем. Поздравляю!"
		}
	}
	return ""
}