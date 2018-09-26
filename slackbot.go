package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	slackbot()
}

func slackbot() {
	port := map[bool]string{true: os.Getenv("PORT"), false: "8080"}[os.Getenv("PORT") != ""]
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	api := slack.New("xoxb-379678856742-414621599025-kscIpN2XZgX6vRNwlkQHT1jC")
	perguntas := []string{"O que você fez ontem?", "O que você fará hoje?", "Existe algum impedimento na sua tarefa?"}
	respostas := make([]string, len(perguntas))
	canal := "infinitydaily"
	acao := " postou no Infinity Daily Standup"
	mensagemFinal := "Bom trabalho!"
	mensagemCancelamento := "Report Cancelado"
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.Text == "cancel" || ev.Text == "Cancel" {
				user, _ := api.GetUserInfo(ev.User)
				_, _, channelID, _ := api.OpenIMChannel(user.ID)
				api.PostMessage(channelID, mensagemCancelamento, slack.PostMessageParameters{})
				slackbot()
			}
			if ev.Text == "report" || ev.Text == "Report" {
				for i := 0; i < len(perguntas); i++ {
					respondeu := false
					user, _ := api.GetUserInfo(ev.User)
					_, _, channelID, _ := api.OpenIMChannel(user.ID)
					api.PostMessage(channelID, perguntas[i], slack.PostMessageParameters{})
					for msg := range rtm.IncomingEvents {
						switch ev := msg.Data.(type) {
						case *slack.MessageEvent:
							if ev.Text == "cancel" || ev.Text == "Cancel" {
								user, _ := api.GetUserInfo(ev.User)
								_, _, channelID, _ := api.OpenIMChannel(user.ID)
								api.PostMessage(channelID, mensagemCancelamento, slack.PostMessageParameters{})
								slackbot()
							}
							if ev.Text != perguntas[i] {
								respostas[i] = ev.Text
								respondeu = true
							}
						}
						if respondeu == true {
							break
						}
					}
				}
				if respostas[len(respostas)-1] != "" {
					user, _ := api.GetUserInfo(ev.User)
					_, _, channelID, _ := api.OpenIMChannel(user.ID)
					api.PostMessage(channelID, mensagemFinal, slack.PostMessageParameters{})
					params := slack.PostMessageParameters{}
					fields := make([]slack.AttachmentField, len(perguntas))
					for i := 0; i < len(perguntas); i++ {
						fields[i].Title = perguntas[i]
						fields[i].Value = respostas[i]
					}
					attachment := slack.Attachment{Fields: fields}
					params.Attachments = []slack.Attachment{attachment}
					params.AsUser = false
					params.IconURL = user.Profile.Image48
					params.Username = user.RealName
					api.PostMessage(canal, user.RealName+acao, params)
					slackbot()
				}
			}
		}
	}
}