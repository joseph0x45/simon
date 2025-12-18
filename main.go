package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func sendMessageToUser(s *discordgo.Session, userID, content string) error {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(channel.ID, content)
	return err
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		currentChannel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Println("[ERROR]: Failed to get channel with ID:", err.Error())
			return
		}
		if currentChannel.Name != "bot" {
			if strings.HasPrefix(m.Content, "m!play") {
				if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
					log.Println("[ERROR]: Failed to delete message with ID:", err.Error())
				}
				mention := fmt.Sprintf("%s", m.Author.Mention())
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hey %s, N'envoie pas les messages la ici singe", mention))
			}
		}
	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /message", func(w http.ResponseWriter, r *http.Request) {
		recipient := r.URL.Query().Get("recipient")
		message := r.URL.Query().Get("message")
		if recipient == "" || message == "" {
			w.WriteHeader(http.StatusOK)
			return
		}
		err = sendMessageToUser(session, recipient, message)
		if err != nil {
			log.Println("[ERROR] Failed to send message to user: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	log.Println("Bot online")
	log.Println("Starting HTTP server on port 6969")
	if err := http.ListenAndServe(":6969", mux); err != nil {
		panic(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-sc
}
