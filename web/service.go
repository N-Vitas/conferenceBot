package web


import (
	"net/http"
	"path"
	"strconv"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"github.com/russross/blackfriday"
	"html/template"
	"conferenceBot/web/realtime"
	"conferenceBot/database"
)

type Post struct {
	Title string
	Body  template.HTML
}

type Service struct {
	Realtime *realtime.ServiceEngine
	Db *database.Connect
}

var (
	// компилируем шаблоны, если не удалось, то выходим
	post_template = template.Must(template.ParseFiles(path.Join("./templates", "layout.html"), path.Join("./templates", "post.html")))
	error_template =  template.Must(template.ParseFiles(path.Join("./templates", "layout.html"), path.Join("./templates", "error.html")))
)

func (self *Service) PostHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	// Извлекаем параметр
	// Например, в http://127.0.0.1:3000/p1 page = "p1"
	// в http://127.0.0.1:3000/ page = ""
	page := params.Get(":page")
	// Путь к файлу (без расширения)
	// Например, view/p1
	p := path.Join("view", page)
	var post_md string
	if page != "" {
		// если page не пусто, то считаем, что запрашивается файл
		// получим view/p1.md
		post_md = p + ".md"
	} else {
		// если page пусто, то выдаем главную
		post_md = p + "/index.html"
	}
	fmt.Println(page)
	// Попытка загрузить файл с расширением md
	post, status, err := load_post(post_md)
	if page == "test" {
		str := ""
		po := params.Get("po")
		dis := params.Get("dis")
		happy := params.Get("happy")
		if po != "" && dis != "" && happy != ""{
			r,e := strconv.ParseFloat(po,64)
			d,e := strconv.ParseFloat(dis,64)
			h,e := strconv.ParseFloat(happy,64)
			fmt.Println("%v",r,d,h,e)
		}
		post_md = p + ".html"
		post, status, err = load_post_html(post_md)
		if str != ""{
			post.Body += template.HTML("<h3 class='fix'>"+str+"</h3>")
		}
	}
	// Если файл с расширением md не найден ищем файл с расширением html
	if err != nil {
		post_md = p + ".html"
		post, status, err = load_post_html(post_md)
	}
	// Если и с расширением html ненайден то выводим ошибку
	if err != nil {
		errorHandler(w, r, status)
		return
	}
	// Собираем общий вид страницы вставляя в нее сформированный файл и выводим на экран. Или показываем ошибку 500
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		fmt.Println(err.Error())
		errorHandler(w, r, 500)
	}
}
// Загружает markdown-файл и конвертирует его в HTML
// Возвращает объект типа Post
// Если путь не существует или является каталогом, то возвращаем ошибку
func load_post(md string) (Post, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не существует
			return Post{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		// не файл, а папка
		return Post{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	fileread, _ := ioutil.ReadFile(md)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	body := strings.Join(lines[1:len(lines)], "\n")
	body = string(blackfriday.MarkdownCommon([]byte(body)))
	post := Post{title, template.HTML(body)}
	return post, 200, nil
}
func load_post_html(md string) (Post, int, error) {
	info, err := os.Stat(md)
	if err != nil {
		if os.IsNotExist(err) {
			// файл не существует
			return Post{}, http.StatusNotFound, err
		}
	}
	if info.IsDir() {
		// не файл, а папка
		return Post{}, http.StatusNotFound, fmt.Errorf("dir")
	}
	fileread, _ := ioutil.ReadFile(md)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	post := Post{title, template.HTML(string(fileread))}
	return post, 200, nil
}
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := error_template.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

