package main

import (
	"fmt"
	"log"
	"net/http"
	"waffe/antibot"
	"waffe/utils"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

func onRequestHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("-------------------\n%s %s\n", r.Method, utils.GetFullPath(r))

		if antibot.NeedsVerification(db, utils.GetIP(r)) {
			if !utils.IsHTMLRequest(r) {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			antibot.StartJavaScriptVerification(db, w, r)
			return
		}

		if !antibot.IsClientVerified(db, utils.GetIP(r)) {
			http.Error(w, "Checks falied", http.StatusForbidden)
			return
		}

		originRequest := utils.ProcessRequest(r)
		utils.RequestOrigin(w, originRequest)
	}
}

func judgeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		antibot.JudgeClient(db, w, r)
	}
}

func main() {
	cfg := utils.LoadConfig("config.yml")
	color.Green("Loaded config file")
	db := antibot.InitDB()
	color.Green("Initialized database")

	if !utils.IsOriginAlive(cfg.Server.Origin) {
		color.Red("Origin server is not reachable! Exiting...")
		color.Red("Please check your origin server address in the config file and make sure it has started.")
		return
	}
	color.Green("Origin server is reachable")

	http.HandleFunc("/.__/api/__judge", judgeHandler(db))
	http.HandleFunc("/", onRequestHandler(db))

	color.Green("Server running at http://%s\n", cfg.Server.Proxy)
	color.Blue("Private IP: %s\n", utils.GetPrivateIP())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(cfg.Server.Proxy), nil))
}
