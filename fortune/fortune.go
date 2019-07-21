package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
)

var fortunes []string

func main() {
	b, err := ioutil.ReadFile("fortunes")
	if err != nil {
		fmt.Println(err)
		return
	}
	fortunes = strings.Split(strings.Replace(string(b), "\r\n", "\n", -1), "\n%\n")
	var keyfile string

	if len(os.Args) < 2 {
		keyfile = "token"
	} else {
		keyfile = os.Args[1]
	}
	keyb, err := ioutil.ReadFile(keyfile)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	key := string(keyb)
	key = strings.TrimSpace(key)
	session, err := discord.New("Bot " + key)
	defer session.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	session.AddHandler(handleMessages)
	_ = session.Open()

	var msg discord.Message

	msg.Content = "Hello world!"
	msg.GuildID = "580751854498218015"
	msg.ChannelID = "580751855030763531"

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func handleMessages(s *discord.Session, m *discord.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) >= 8 && m.Content[:8] == "!fortune" {
		s.ChannelMessageSend(m.ChannelID, fortunes[rand.Intn(len(fortunes))])
	}
}
