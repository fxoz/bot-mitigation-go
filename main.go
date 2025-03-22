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

		if !antibot.IsClientWhitelisted(db, utils.GetIP(r)) {
			if !utils.IsHTMLRequest(r) {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			fmt.Println("Protecting path:", r.URL.Path)
			antibot.BotProtectionHandler(db, w, r)
			return
		}

		originRequest := utils.ProcessRequest(r)
		utils.RequestOrigin(w, originRequest)
	}
}

func verifyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("-------------------\n%s %s\n", r.Method, utils.GetFullPath(r))
		antibot.VerifyHandler(db, w, r)
	}
}

func main() {
	cfg := utils.LoadConfig("config.yml")
	db := antibot.InitDB()
	color.Green("Initialized database")

	if !utils.IsOriginAlive(cfg.Server.Origin) {
		color.Red("Origin server is not reachable! Exiting...")
		color.Red("Please check your origin server address in the config file and make sure it has started.")
		return
	}
	color.Green("Origin server is reachable")

	http.HandleFunc("/__verify", verifyHandler(db))
	http.HandleFunc("/__verify/", verifyHandler(db))
	http.HandleFunc("/", onRequestHandler(db))

	color.Green("Server running at http://%s\n", cfg.Server.Proxy)
	color.Blue("Private IP: %s\n", utils.GetPrivateIP())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(cfg.Server.Proxy), nil))
}
