package main

import (
	"./api"
	// "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	// "encoding/hex"
	"net/http"
	// "strconv"
	"crypto/sha1"
	"strings"
)

func main() {
	fmt.Println("come on")

	http.HandleFunc("/", HttpWrapHandler(get))
	http.HandleFunc("/add", HttpWrapHandler(add))
	http.HandleFunc("/del", HttpWrapHandler(del))
	http.HandleFunc("/edit", HttpWrapHandler(edit))
	http.HandleFunc("/login", HttpWrapHandler(login))
	http.HandleFunc("/out", HttpWrapHandler(out))
	http.HandleFunc("/register", HttpWrapHandler(register))

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

type HttpErrorHandler func(w http.ResponseWriter, r *http.Request) error
type HttpHandler func(w http.ResponseWriter, r *http.Request)

func HttpWrapHandler(inHandler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				if err == api.NoLoginError {
					http.Redirect(w, r, "/login", 302)
				} else {
					fmt.Println(err)
					fmt.Fprint(w, "<div style=\"color:red\">")
					fmt.Fprintf(w, "%v", err)
					fmt.Fprint(w, "</div>")
				}
			}
			// panic(err)
		}()
		inHandler(w, r)

		fmt.Println("come out")
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	v := api.Get("SELECT * FROM book")

	result := api.TemplateOutput("index.html", struct {
		List    []api.Book
		IsLogin bool
	}{
		v,
		true,
	})

	w.Write(result)
}

func add(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	if r.Method == "GET" {
		result := api.TemplateOutput("add.html", nil)

		w.Write(result)
	} else {
		api.CheckCsrf(w, r)

		data := api.CheckInput(r, map[string]string{
			"username": "string",
			"bname":    "string",
		})

		_ = api.Add(
			"INSERT book SET Username=?,Bname=?",
			data["username"].(string),
			data["bname"].(string),
		)

		http.Redirect(w, r, "/", 302)
	}

}

func del(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	data := api.CheckInput(r, map[string]string{
		"id": "int",
	})

	v := api.Get("SELECT * FROM book where Bid=?", data["id"].(int))

	if len(v) == 0 {
		panic(errors.New("不存在该数据"))
	}

	api.Del(data["id"].(int))

	http.Redirect(w, r, "/", 302)

}

func edit(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	if r.Method == "GET" {
		data := api.CheckInput(r, map[string]string{
			"id": "int",
		})

		v := api.Get("SELECT * FROM book where Bid=?", data["id"].(int))

		if len(v) == 0 {
			panic(errors.New("不存在该数据"))
		}

		result := api.TemplateOutput("edit.html", v[0])

		w.Write(result)
	} else {
		api.CheckCsrf(w, r)
		data := api.CheckInput(r, map[string]string{
			"id":    "int",
			"bname": "string",
		})

		api.Edit(
			"update book set Bname=? where Bid=?",
			data["bname"].(string),
			data["id"].(int),
		)

		http.Redirect(w, r, "/", 302)
	}

}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		result := api.TemplateOutput("login.html", nil)

		w.Write(result)
	} else {

		data := api.CheckInput(r, map[string]string{
			"username": "string",
			"password": "string",
		})

		v := api.GetUserinfo("SELECT * FROM userinfo where username=?", data["username"])

		if len(v) == 0 {
			panic(errors.New("用户不存在"))
		}

		fmt.Println(v)

		hashAndSalt := strings.Split(v[0].Password, ":")
		password := hashAndSalt[0]
		salt := hashAndSalt[1]
		hash := sha1.New()
		passwordSha1Byte := hash.Sum([]byte(data["password"].(string) + salt))
		passwordSha1 := hex.EncodeToString(passwordSha1Byte)

		if password != passwordSha1 {
			fmt.Println("数据库password:", v[0].Password)
			fmt.Println("生成出来的密码:", passwordSha1)
			fmt.Println("数据库密码:", password)
			fmt.Println("盐:", salt)

			panic(errors.New("密码错误"))
		}

		ss := api.SessionInit(w, r)
		defer ss.SessionClose()
		// ss.SessionSet("username", "admin")
		ss.SessionSet("username", v[0].Username)

		http.Redirect(w, r, "/", 302)
	}

}

func out(w http.ResponseWriter, r *http.Request) {

	s := api.SessionInit(w, r)
	s.SessionDel("username")
	defer s.SessionClose()

	http.Redirect(w, r, "/", 302)

}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		result := api.TemplateOutput("register.html", nil)

		w.Write(result)
	} else {
		data := api.CheckInput(r, map[string]string{
			"username":  "string",
			"password":  "string",
			"password2": "string",
		})

		if data["password"] != data["password2"] {
			panic(errors.New("确认密码不正确"))
		}

		username := data["username"].(string)
		password := data["password"].(string)

		v := api.GetUserinfo("SELECT * FROM userinfo where username=?", username)

		if len(v) > 0 {
			panic(errors.New("用户名已存在，请重新注册其他用户名字"))
		}

		salt := api.RandString(5)
		fmt.Println("salt", salt)
		hash := sha1.New()
		passwordSha1Byte := hash.Sum([]byte(password + salt))
		passwordSha1 := hex.EncodeToString(passwordSha1Byte)

		fmt.Println("passwordSha1", passwordSha1)

		_ = api.Add(
			"INSERT userinfo SET username=?,password=?",
			username,
			passwordSha1+":"+salt,
		)

		ss := api.SessionInit(w, r)
		defer ss.SessionClose()
		ss.SessionSet("username", username)

		http.Redirect(w, r, "/", 302)
	}

}
