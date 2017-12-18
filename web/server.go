package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-dashboard/analyzer"
	"github.com/zhangmingkai4315/dns-dashboard/model"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

var store *sessions.CookieStore

// Message 定义基本的response json响应对象
type Message struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// auth 中间件提供认证页面的授权和跳转
func auth(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "dns-dashboard")
		user, ok := session.Values["username"].(string)
		if !ok || user == "" {
			http.Redirect(w, req, "/login", 302)
			return
		}
		context.Set(req, "username", user)
		f(w, req)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	data := map[string]interface{}{
		"username": context.Get(req, "username"),
	}
	r.HTML(w, http.StatusOK, "index", data)
}

func login(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "dns-dashboard")
	r := render.New(render.Options{})
	if user, ok := session.Values["username"].(string); ok && user != "" {
		http.Redirect(w, req, "/", 302)
		return
	}
	config, err := utils.GetConfig()
	if err != nil {
		r.HTML(w, http.StatusServiceUnavailable, "login", nil)
		return
	}
	if req.Method == "POST" {
		user := req.FormValue("username")
		password := req.FormValue("password")
		if user == config.Username && password == config.Password {
			session.Values["username"] = user
			session.Save(req, w)
			http.Redirect(w, req, "/", 302)
		} else {
			http.Redirect(w, req, "/login", 401)
		}
	} else if req.Method == "GET" {
		r := render.New(render.Options{})
		r.HTML(w, http.StatusOK, "login", nil)
	}
}

// 系统的logout处理函数
func logout(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "dns-dashboard")
	if user, ok := session.Values["username"].(string); ok && user != "" {
		session.Values["username"] = ""
		session.Save(req, w)
		http.Redirect(w, req, "/login", 302)
		return
	}
}

// getStatus 定义api:/status路由函数，返回系统的实时信息
func getStatus(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	status := analyzer.GetSystemStatus()
	r.JSON(w, http.StatusOK, Message{Data: status})
}

// getDNSStatus 定义api:/dns_init_status路由函数,返回dns的信息
func getDNSStatus(w http.ResponseWriter, req *http.Request) {
	db, err := model.GetDB()
	r := render.New(render.Options{})
	if err != nil {
		r.JSON(w, http.StatusServiceUnavailable, Message{Error: err.Error()})
		return
	}
	var serials []model.DNSSerialData
	if err := db.Order("timestamp desc").Limit(10).Find(&serials).Error; err != nil {
		r.JSON(w, http.StatusServiceUnavailable, Message{Error: err.Error()})
	} else {
		r.JSON(w, http.StatusOK, Message{Data: serials})
	}
}

// getLastestDNSStatus 定义api接口函数，返回最新的一条dns的状态信息
func getLastestDNSStatus(w http.ResponseWriter, req *http.Request) {
	db, err := model.GetDB()
	r := render.New(render.Options{})
	if err != nil {
		r.JSON(w, http.StatusServiceUnavailable, Message{Error: err.Error()})
		return
	}
	var serials model.DNSSerialData
	if err := db.Order("timestamp desc").Limit(1).Find(&serials).Error; err != nil {
		r.JSON(w, http.StatusServiceUnavailable, Message{Error: err.Error()})
	} else {
		r.JSON(w, http.StatusOK, Message{Data: serials})
	}
}

// StartServer 启动web服务器并定义所有路由接口API
func StartServer(config *utils.Config) {
	r := mux.NewRouter()
	store = sessions.NewCookieStore([]byte(config.Global.Secrect))
	r.HandleFunc("/", auth(index)).Methods("GET")
	r.HandleFunc("/logout", logout).Methods("POST", "GET")
	r.HandleFunc("/login", login).Methods("GET", "POST")

	r.HandleFunc("/status", getStatus).Methods("GET")
	r.HandleFunc("/dns_init_status", getDNSStatus).Methods("GET")
	r.HandleFunc("/dns_lastest_status", getLastestDNSStatus).Methods("GET")
	serverWithPort := fmt.Sprintf("%s:%d", config.Global.Server, config.Global.Port)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	log.Printf("Server listen : %s", serverWithPort)
	log.Panic(http.ListenAndServe(serverWithPort, http.TimeoutHandler(r, time.Second*10, "timeout")))
}
