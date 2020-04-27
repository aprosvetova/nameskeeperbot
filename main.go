package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

var db *redis.Client
var bot *tgbotapi.BotAPI
var cfg *Config

func main() {
	var err error

	cfg, err = loadConfig()
	if err != nil {
		log.Fatalln("can't read env", err)
	}

	db = redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	bot, err = tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatalln("can't access Bot API: ", err)
	}

	if cfg.TdlibEnabled {
		go listenTdlib()
	}

	log.Printf("Started listening on @%s\n", bot.Self.UserName)

	updates, err := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message
		go saveNameWithBot(msg.From)
		if msg.Chat.Type == "private" {
			if msg.Command() == "start" {
				handleStart(msg)
				continue
			}
			if msg.ForwardFrom != nil {
				handleSearch(msg, msg.ForwardFrom.ID)
				go saveNameWithBot(msg.ForwardFrom)
				continue
			}
			handleUsage(msg)
			continue
		}
		if msg.Command() == "names" {
			if msg.ReplyToMessage == nil {
				handleUsage(msg)
				continue
			}
			handleSearch(msg, msg.ReplyToMessage.From.ID)
			go saveNameWithBot(msg.ReplyToMessage.From)
		}
	}
}

func handleStart(msg *tgbotapi.Message) {
	_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Hey there! I'm Names Keeper.\n"+
		"I can show you one's name history.\n\n"+
		"There are two ways to ask me for that:\n"+
		"1. Forward me one's message privately\n"+
		"2. Reply to one's message with /names command while being in group\n\n"+
		"Please note that I learn names listening to groups so if I don't know one's history, he/she hasn't been chatting in group where I exist while changing names."))
}

func handleUsage(msg *tgbotapi.Message) {
	c := tgbotapi.NewMessage(msg.Chat.ID, "Hey, I do only work with forwards in private and with /names command replied on target user in groups")
	c.ReplyToMessageID = msg.MessageID
	_, _ = bot.Send(c)
}

func handleSearch(replyTo *tgbotapi.Message, targetID int) {
	msg := getNamesMessage(targetID)
	c := tgbotapi.NewMessage(replyTo.Chat.ID, msg)
	c.ReplyToMessageID = replyTo.MessageID
	_, _ = bot.Send(c)
}

func saveNameWithBot(u *tgbotapi.User) {
	saveName(u.ID, u.FirstName, u.LastName, u.UserName)
}

func saveName(userID int, firstName, lastName, username string) {
	latestName := getLatestName(userID)

	currentName := strings.TrimSpace(firstName + " " + lastName)
	if username != "" {
		currentName += " @" + username
	}

	if latestName != "" && latestName != currentName {
		lastChanged := getSetLastChanged(userID, time.Now())
		if lastChanged != nil && time.Since(*lastChanged) < 5*time.Minute {
			db.ZRem(getUserKey(userID), latestName)
		}
	}

	storeName(userID, currentName)
}

func getNamesMessage(userID int) (message string) {
	records := getNames(userID)
	if len(records) == 0 {
		return "I haven't learned any names of this user :(\nTry adding me to the group where he/she talks frequently."
	}
	for i, record := range records {
		lastSeen := "Last known"
		if i != 0 {
			lastSeen = "Until " + record.LastSeen.Format("02.01.2006")
		}
		message += fmt.Sprintf("%s: %s\n", lastSeen, record.Name)
	}
	return
}
