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
	ontem, hoje, impedimento := "", "", ""
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if len(ev.User) == 0 {
				continue
			}
			if ev.Text == "cancel" || ev.Text == "Cancel" {
				user, err := api.GetUserInfo(ev.User)
				_, _, channelID, err := api.OpenIMChannel(user.ID)
				if err != nil {
					fmt.Printf("[WARN]  could not grab user information: %s", ev.User)
				}
				api.PostMessage(channelID, "Report cancelado", slack.PostMessageParameters{})
			}
			if ev.Text == "report" || ev.Text == "Report" {
				user, err := api.GetUserInfo(ev.User)
				_, _, channelID, err := api.OpenIMChannel(user.ID)
				if err != nil {
					fmt.Printf("[WARN]  could not grab user information: %s", ev.User)
				}
				api.PostMessage(channelID, "O que você fez ontem?", slack.PostMessageParameters{})
				for msg := range rtm.IncomingEvents {
					switch ev := msg.Data.(type) {
					case *slack.MessageEvent:
						if ev.Text == "cancel" || ev.Text == "Cancel" {
							user, err := api.GetUserInfo(ev.User)
							_, _, channelID, err := api.OpenIMChannel(user.ID)
							if err != nil {
								fmt.Printf("[WARN]  could not grab user information: %s", ev.User)
							}
							api.PostMessage(channelID, "Report cancelado", slack.PostMessageParameters{})
							slackbot()
						}
						if ev.Text != "O que você fez ontem?" && ev.Text != "cancel" {
							ontem = ev.Text
						}
						if len(ontem) > 0 {
							api.PostMessage(channelID, "O que você fará hoje?", slack.PostMessageParameters{})
							for msg := range rtm.IncomingEvents {
								switch ev := msg.Data.(type) {
								case *slack.MessageEvent:
									if ev.Text == "cancel" || ev.Text == "Cancel" {
										user, err := api.GetUserInfo(ev.User)
										_, _, channelID, err := api.OpenIMChannel(user.ID)
										if err != nil {
											fmt.Printf("[WARN]  could not grab user information: %s", ev.User)
										}
										api.PostMessage(channelID, "Report cancelado", slack.PostMessageParameters{})
										slackbot()
									}
									if ev.Text != "O que você fará hoje?" && ev.Text != "cancel" {
										hoje = ev.Text
									}
									if len(hoje) > 0 {
										api.PostMessage(channelID, "Existe algum impedimento na sua tarefa?", slack.PostMessageParameters{})
										for msg := range rtm.IncomingEvents {
											switch ev := msg.Data.(type) {
											case *slack.MessageEvent:
												if ev.Text == "cancel" || ev.Text == "Cancel" {
													user, err := api.GetUserInfo(ev.User)
													_, _, channelID, err := api.OpenIMChannel(user.ID)
													if err != nil {
														fmt.Printf("[WARN]  could not grab user information: %s", ev.User)
													}
													api.PostMessage(channelID, "Report cancelado", slack.PostMessageParameters{})
													slackbot()
												}
												if ev.Text != "Existe algum impedimento na sua tarefa?" && ev.Text != "cancel" {
													impedimento = ev.Text
													if len(impedimento) > 0 {
														api.PostMessage(channelID, "Bom trabalho!", slack.PostMessageParameters{})
														params := slack.PostMessageParameters{}
														attachment := slack.Attachment{

															Fields: []slack.AttachmentField{
																slack.AttachmentField{
																	Title: "O que você fez ontem?",
																	Value: ontem,
																},
																slack.AttachmentField{
																	Title: "O que você fará hoje?",
																	Value: hoje,
																},
																slack.AttachmentField{
																	Title: "Existe algum impedimento na sua tarefa?",
																	Value: impedimento,
																},
															},
														}
														params.Attachments = []slack.Attachment{attachment}
														params.AsUser = false
														params.IconURL = user.Profile.Image48
														params.Username = user.RealName
														api.PostMessage("infinitydaily", user.RealName+" postou no Infinity Daily Standup", params)
														slackbot()
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

		case *slack.RTMError:
			fmt.Printf("[ERROR] %s\n", ev.Error())
		}
	}
}
