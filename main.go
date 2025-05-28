package main

import (
	"chocolate/internal/csrfKey"
	"chocolate/internal/passwords"
	"chocolate/pkg/commands"
	"chocolate/pkg/db"

	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const frontEndFilesDir string = "./web"
const htmlDirPath string = "./web/html"
const imgFileExt string = ".png"
const imgDirPath string = "./web/img"
const dbName string = "../data/chocolate.db"
const frontDirName string = "/web/"

var defaultCategories = []string{"Наборы", "Плитки шоколада", "Батончики"}

func main() {
	fmt.Println(`
Запуск программы

Справка по пользованию:
s - остановка программы
filldb - заполнение базовых категорий и админов
additems - добавить несколько товаров`)
	db.OpenDB(dbName)

	db.FillDefaultData(defaultCategories)

	go WaitForCommands()

	checkFlags()
	router := mux.NewRouter()
	router.PathPrefix(frontDirName).Handler(http.StripPrefix(frontDirName,
		http.FileServer(http.Dir(frontEndFilesDir))))
	//GET запросы
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/admin", adminPageHandler)
	router.HandleFunc("/product", productPageHandler)
	router.HandleFunc("/login", loginAdminPageHandler)
	router.HandleFunc("/search_product", searchProductHandler)
	router.HandleFunc("/logout", logoutAdminPageHandler)
	//POST запросы
	router.HandleFunc("/add_product", addProductHandler)
	router.HandleFunc("/add_admin", addAdminHandler)
	router.HandleFunc("/add_category", addCategoryHandler)
	//DELETE запросы
	router.HandleFunc("/delete_category", deleteCategoryHandler)
	router.HandleFunc("/delete_product", deleteProductHandler)

	fmt.Println("Сервер будет запущен на http://0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", csrf.Protect([]byte(csrfKey.CSRF_KEY), csrf.Secure(false))(router)) // поставить на тру при девелопе
	if err != nil {
		panic(err)
	}
}
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path.Join(htmlDirPath, "/index.html"))
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "index", prepareForIndex())
}

type IndexCategory struct {
	Name  string
	Items []db.Item
}
type IndexData struct {
	Categories []IndexCategory
}

func prepareForIndex() IndexData {
	var data IndexData
	categories := db.FindAllCategories()
	for _, category := range categories {
		data.Categories = append(data.Categories, IndexCategory{Name: category, Items: db.FindCategoryItems(category)})

	}
	return data
}
func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join(htmlDirPath, "admin.html"))
	if err != nil {
		panic(err)
	}
	allCategories := db.FindAllCategories()
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "admin", map[string]any{csrf.TemplateTag: csrf.TemplateField(r), "allCategories": allCategories})
}
func loginAdminPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		name, password := r.Form.Get("name"), r.Form.Get("password")
		adminsPasswords := db.FindAdminsPasswords(name)
		for _, hashPassword := range adminsPasswords {
			if passwords.CheckPassword(password, hashPassword) {
				http.SetCookie(w, &http.Cookie{
					Name:  "isAuthenticated",
					Value: "true",
				})
				http.Redirect(w, r, "/admin", http.StatusFound)
			}
		}
	}
	tmpl, err := template.ParseFiles(path.Join(htmlDirPath, "login.html"))
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "login", map[string]any{csrf.TemplateTag: csrf.TemplateField(r)})
}
func logoutAdminPageHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("ok"))
}
func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("image")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Сохраняем файл на диск
		r.ParseForm()
		name, price_str, long_description, category, short_description := r.Form.Get("name"), r.Form.Get("price"), r.Form.Get("description"),
			r.Form.Get("category"), r.Form.Get("short_description")
		price, _ := strconv.Atoi(price_str)
		fileName := path.Join(imgDirPath, name+imgFileExt)
		f, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			panic(err)
		}

		db.CreateItem(name, long_description, short_description, uint(price), db.Category{Name: category})
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func addAdminHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name, password := r.Form.Get("name"), r.Form.Get("password")

	db.CreateAdmin(name, password, true)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

type deleteProductJson struct {
	Product  string `json:"product"`
	Category string `json:"category"`
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {

		defer r.Body.Close()
		jsonData := deleteProductJson{}
		json.NewDecoder(r.Body).Decode(&jsonData)
		name, category := jsonData.Product, jsonData.Category

		db.DeleteItem(name, category)
		os.Remove(path.Join(imgDirPath, name+imgFileExt))
		fmt.Println("Продукт", name, "удален из категории", category)
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func searchProductHandler(w http.ResponseWriter, r *http.Request) {
	search_name := r.URL.Query().Get("name")

	type Item struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	items := []Item{}

	for _, item := range db.FindItems(search_name) {
		items = append(items, Item{Name: item.Name, Category: item.Category.Name})
	}
	json.NewEncoder(w).Encode(items)
}
func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		r.ParseForm()

		name := r.Form.Get("name")

		db.CreateCategory(name)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

type deleteCategoryJson struct {
	Category string `json:"name"`
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		defer r.Body.Close()
		jsonData := deleteCategoryJson{}
		json.NewDecoder(r.Body).Decode(&jsonData)

		db.DeleteCategory(jsonData.Category)
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func productPageHandler(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query()["name"][0]
	tmpl, err := template.ParseFiles(path.Join(htmlDirPath, "product.html"))
	if err != nil {
		panic(err)
	}
	item := db.FindItems(productName)[0]
	//выводим шаблон клиенту в браузер
	err = tmpl.ExecuteTemplate(w, "product", item)
	if err != nil {
		panic(err)
	}
}
func WaitForCommands() {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "s":
			commands.Stop()
		case "filldb":
			db.FillDefaultData(defaultCategories)
		case "additems":
			db.AddItems(defaultCategories)
		}
	}
}
func checkFlags() {
	additems := flag.Bool("additems", false, "добавить несколько товаров")
	flag.Parse()
	if *additems {
		fmt.Println("Выполняю переданный флаг additems:")
		db.AddItems(defaultCategories)
	}
}
