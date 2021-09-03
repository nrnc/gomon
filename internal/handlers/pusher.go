package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/pusher/pusher-http-go"
)

func (repo *DBRepo) PusherAuth(rw http.ResponseWriter, r *http.Request) {
	userId := repo.App.Session.GetInt(r.Context(), "userID")

	u, _ := repo.DB.GetUserById(userId)
	params, _ := ioutil.ReadAll(r.Body)
	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(userId),
		UserInfo: map[string]string{
			"name": u.FirstName,
			"id":   strconv.Itoa(userId),
		},
	}
	response, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
	if err != nil {
		log.Println(err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(response)
}

func (repo *DBRepo) PusherTest(rw http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["message"] = "Hello World!"
	err := repo.App.WsClient.Trigger("public-channel", "test-event", data)
	if err != nil {
		log.Println(err)
		return
	}
}
