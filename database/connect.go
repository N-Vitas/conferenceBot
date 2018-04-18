package database

import (
	"database/sql"
	. "conferenceBot/config"
	_ "github.com/mattn/go-sqlite3"
	"unicode"
	"os"
	"fmt"
)

// Сессия с файловой базой данных
type Connect struct {
	db *sql.DB
}
func NewConnect() *Connect{
	db := &Connect{}
	f, err := os.Open(DATABASE_FILE)
	if err != nil {
		f, err = os.Create(DATABASE_FILE)
		if err != nil {
			panic("[database][error] Нет файла базы данных")
		}
	}
	f.Close();
	db.migrate()
	return db
}
func (self *Connect) migrate()  {
	query := []string{
		"CREATE TABLE IF NOT EXISTS users (chat_id NUMERIC, phone TEXT, name TEXT, login TEXT, photo TEXT, status NUMERIC);",
		"CREATE TABLE IF NOT EXISTS queue (id INTEGER PRIMARY KEY, question TEXT, option1 TEXT, option2 TEXT, option3 TEXT, option4 TEXT, answer TEXT);",
		"CREATE TABLE IF NOT EXISTS games (id INTEGER PRIMARY KEY, chat_id NUMERIC, countErr NUMERIC, countRight NUMERIC, done TEXT, status NUMERIC);",
	}
	for _,v := range query{
		smtp, err := self.GetDb().Prepare(v)
		defer smtp.Close()
		self.checkErr(err)
		_, err = smtp.Exec()
		self.checkErr(err)
	}
}
// Создание новой записи
func (self *Connect) CreateUser(user User) int64 {
	smtp,err := self.GetDb().Prepare("INSERT INTO users( chat_id,phone,[name],login,photo,status ) VALUES( ?, ?, ?, ?, ?, ? )")
	self.checkErr(err)
	status := 0
	if user.Status {
		status = 1
	}
	res,err := smtp.Exec(user.Chat_ID,user.Phone,user.Name,user.Login,user.Photo,status)
	id,err := res.LastInsertId()
	self.checkErr(err)
	smtp.Close()
	return id
}

// Обновление записи пользователя
func (self *Connect) UpdateUser(user User) int64 {
	smtp,err := self.GetDb().Prepare("UPDATE users SET phone=?, name=?,status=?,login=?,photo=? WHERE chat_id=?")
	self.checkErr(err)
	status := 0
	if user.Status {
		status = 1
	}
	res,err := smtp.Exec(user.Phone,user.Name,status,user.Login,user.Photo,user.Chat_ID)
	id,err := res.LastInsertId()
	self.checkErr(err)
	smtp.Close()
	return id
}
// Получение списка пользователей
func (self *Connect) GetUsers() map[int64]User {
	var (
		chat_id sql.NullInt64
		phone sql.NullString
		name sql.NullString
		login sql.NullString
		photo sql.NullString
		status sql.NullInt64
	)
	Result := map[int64]User{}
	rows,err := self.GetDb().Query("SELECT chat_id,phone,name,login,photo,status FROM users")
	self.checkErr(err)
	for rows.Next() {
		err := rows.Scan(&chat_id,&phone,&name,&login,&photo,&status)
		if err != nil {
			continue
		}
		if chat_id.Valid {
			uStatus := false
			if status.Int64 == 1 {
				uStatus = true
			}
			Result[chat_id.Int64] = User{
				Chat_ID:chat_id.Int64,
				Phone:phone.String,
				Name:name.String,
				Login:login.String,
				Photo:photo.String,
				Status: uStatus,
			}
		}
	}
	rows.Close()
	return Result
}

// Получение списка пользователей
func (self *Connect) FindUser(chatId int64) (User,error) {
	var (
		chat_id sql.NullInt64
		phone sql.NullString
		name sql.NullString
		login sql.NullString
		photo sql.NullString
		status sql.NullInt64
	)
	Result := User{}
	err := self.GetDb().QueryRow("SELECT chat_id,phone,name,login,photo,status FROM users").Scan(&chat_id,&phone,&name,&login,&photo,&status)
	if self.checkErr(err) {
		return Result,err
	}
	uStatus := false
	if status.Int64 == 1 {
		uStatus = true
	}
	Result = User{
		Chat_ID:chat_id.Int64,
		Phone:phone.String,
		Name:name.String,
		Login:login.String,
		Photo:photo.String,
		Status: uStatus,
	}
	return Result,nil
}
// Переподключение или создание подключения к базе
func (self *Connect) GetDb() *sql.DB {
	if self.db != nil {
		return self.db
	}

	db, err := sql.Open("sqlite3", DATABASE_FILE)
	self.checkErr(err)
	self.db = db
	return self.db
}
// Создание новой записи
func (self *Connect) CreateGameQuestion(q GameQuestion) bool {
	id := int64(0)
	self.GetDb().QueryRow("SELECT MAX(ID) FROM queue").Scan(&id)
	smtp,err := self.GetDb().Prepare("INSERT INTO queue ( id, question, option1, option2, option3, option4, answer ) VALUES( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {	return false }
	res,err := smtp.Exec(id+1,q.Question,q.Option1,q.Option2,q.Option3,q.Option4,q.Answer)
	id,err = res.LastInsertId()
	if err != nil {	return false }
	smtp.Close()
	return true
}
// Создание новой записи
func (self *Connect) GetGameQuestions() []GameQuestion {
	var (
		id sql.NullInt64
		question sql.NullString
		option1 sql.NullString
		option2 sql.NullString
		option3 sql.NullString
		option4 sql.NullString
		answer sql.NullString
	)
	Result := []GameQuestion{}
	rows,err := self.GetDb().Query("SELECT id, question, option1, option2, option3, option4, answer FROM queue")
	self.checkErr(err)
	for rows.Next() {
		err := rows.Scan(&id,&question,&option1,&option2,&option3,&option4,&answer)
		if err != nil {
			continue
		}
		Result = append(Result,GameQuestion{
			id.Int64,
			question.String,
			option1.String,
			option2.String,
			option3.String,
			option4.String,
			answer.String,
		})
	}
	rows.Close()
	return Result
}

// Создание новой записи
func (self *Connect) CreateGamePlayer(q GamePlayUser) bool {
	id := int64(0)
	err := self.GetDb().QueryRow("SELECT id FROM games WHERE chat_id=?",q.Chat_id).Scan(&id)
	if err == sql.ErrNoRows {
		self.GetDb().QueryRow("SELECT MAX(ID) FROM games").Scan(&id)
		smtp,err := self.GetDb().Prepare("INSERT INTO games ( id, chat_id, countErr, countRight, done, status ) VALUES( ?, ?, ?, ?, ?, ? )")
		if err != nil {	return false }
		_,err = smtp.Exec(id+1,q.Chat_id,q.CountErr,q.CountRight,q.Done,q.Status)
		if err != nil {	return false }
		smtp.Close()
	}
	return true
}
// Обновление записи пользователя
func (self *Connect) UpdateGamePlayer(user GamePlayUser) bool {
	smtp,err := self.GetDb().Prepare("UPDATE games SET countErr=?, countRight=?, done=?, status=? WHERE chat_id=?")
	self.checkErr(err)
	if err != nil {	return false }
	_,err = smtp.Exec(user.CountErr,user.CountRight,user.Done,user.Status,user.Chat_id)
	self.checkErr(err)
	if err != nil {	return false }
	smtp.Close()
	return true
}
// Закрытие сессии
func (self *Connect) Close()  {
	if self.GetDb() != nil {
		self.GetDb().Close()
	}
}
// Обработка ошибки
func (self *Connect) checkErr(err error) bool {
	if err != nil {
		return true
		fmt.Println("[database][error] ",err.Error())
	}
	return false
}

func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
