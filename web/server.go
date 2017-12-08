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

type Message struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func getStatus(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	status := analyzer.GetSystemStatus()
	r.JSON(w, http.StatusOK, status)
}

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

// StartServer will start analyzer and the web serve
func StartServer(config *utils.Config) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	}).Methods("GET")

	r.HandleFunc("/status", getStatus).Methods("GET")
	r.HandleFunc("/dns_status", getDNSStatus).Methods("GET")

	serverWithPort := fmt.Sprintf("%s:%d", config.Server, config.Port)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	log.Printf("Server listen : %s", serverWithPort)
	http.ListenAndServe(serverWithPort, http.TimeoutHandler(r, time.Second*10, "timeout"))
}
