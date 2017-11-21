package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/zhangmingkai4315/dns-dashboard/utils"
)

func getStatus(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	r.JSON(w, http.StatusOK, map[string]string{"status": "no update"})
}

// StartServer will start analyzer and the web serve
func StartServer(config *utils.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	}).Methods("GET")
	r.HandleFunc("/status", getStatus).Methods("GET")

	serverWithPort := fmt.Sprintf("%s:%d", config.Server, config.Port)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./web/assets"))))
	log.Printf("Server listen : %s", serverWithPort)
	http.ListenAndServe(serverWithPort, http.TimeoutHandler(r, time.Second*10, "timeout"))
}
