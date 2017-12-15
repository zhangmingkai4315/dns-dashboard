package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-dashboard/analyzer"
	"github.com/zhangmingkai4315/dns-dashboard/model"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

// Message 定义基本的response json响应对象
type Message struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// getStatus 定义api:/status路由函数，返回系统的实时信息
func getStatus(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	status := analyzer.GetSystemStatus()
	r.JSON(w, http.StatusOK, status)
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
		r.JSON(w, http.StatusOK, serials)
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
		r.JSON(w, http.StatusOK, serials)
	}
}

// StartServer 启动web服务器并定义所有路由接口API
func StartServer(config *utils.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	}).Methods("GET")
	r.HandleFunc("/status", getStatus).Methods("GET")
	r.HandleFunc("/dns_init_status", getDNSStatus).Methods("GET")
	r.HandleFunc("/dns_lastest_status", getLastestDNSStatus).Methods("GET")
	serverWithPort := fmt.Sprintf("%s:%d", config.Global.Server, config.Global.Port)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	log.Printf("Server listen : %s", serverWithPort)
	http.ListenAndServe(serverWithPort, http.TimeoutHandler(r, time.Second*10, "timeout"))
}
