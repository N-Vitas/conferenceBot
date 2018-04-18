package config

import (
	"strings"
	"strconv"
)

type Emit struct {
	Action string `json:"action"`
	Body Content `json:"body"`
	User User `json:"user"`
}
type Content struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}
type User struct {
	Chat_ID int64 `json:"chat_id"`
	Phone string `json:"phone"`
	Name string `json:"name"`
	Login string `json:"login"`
	Photo string `json:"photo"`
	Status bool `json:"status"`
}
type GamePlayUser struct {
	Id int64 `json:"id"`
	Chat_id int64 `json:"chat_id"`
	CountErr int64 `json:"count_err"`
	CountRight int64 `json:"count_right"`
	Status int64 `json:"status"`
	Done string `json:"done"`
}
type FullUser struct {
	Chat_ID int64 `json:"chat_id"`
	Phone string `json:"phone"`
	Name string `json:"name"`
	Login string `json:"login"`
	Photo string `json:"photo"`
	Id int64 `json:"id"`
	CountErr int64 `json:"count_err"`
	CountRight int64 `json:"count_right"`
	Status int64 `json:"status"`
	Done string `json:"done"`
}
type GameQuestion struct {
	Id int64 `json:"id"`
	Question string `json:"question"`
	Option1  string `json:"option_1"`
	Option2  string `json:"option_2"`
	Option3  string `json:"option_3"`
	Option4  string `json:"option_4"`
	Answer  string `json:"answer"`
}
/*
 * Вспомогательные функции трансформации
 */
func (u *GamePlayUser) GetAnswersId() []string{
	return strings.Fields(u.Done)
}
func (u *GamePlayUser) AddAnswerId(id int64){
	res := strings.Fields(u.Done)
	res = append(res,strconv.Itoa(int(id)))
	u.Done = strings.Join(res,",")
}
func (u *GamePlayUser) ClearAnswerId(){
	res := []string{}
	u.CountRight = 0
	u.Done = strings.Join(res,",")
}
func (u *GamePlayUser) UpCountRight(){
	u.CountRight = u.CountRight + 1
}
func (u *GamePlayUser) UpCountErr(){
	u.CountErr = u.CountErr + 1
}