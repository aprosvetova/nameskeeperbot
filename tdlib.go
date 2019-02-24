package main

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

func listenTdlib() {
	client.SetLogVerbosityLevel(0)

	authorizer := client.BotAuthorizer(cfg.General.Token)
	authorizer.TdlibParameters <- &client.TdlibParameters{
		DatabaseDirectory:      filepath.Join("tdlib", "database"),
		UseChatInfoDatabase:    true,
		ApiId:                  int32(cfg.TdLib.ApiID),
		ApiHash:                cfg.TdLib.ApiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Ebik",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalln("tdlib error", err)
	}

	_, err = tdlibClient.GetMe()
	if err != nil {
		log.Fatalln("tdlib error", err)
	}

	log.Println("Started listening with tdlib")

	listener := tdlibClient.GetListener()
	defer listener.Close()

	for update := range listener.Updates {
		if update.GetClass() == client.ClassUpdate {
			switch u := update.(type) {
			case *client.UpdateUser:
				saveNameWithTdLib(u.User)
			}
		}
	}
}

func saveNameWithTdLib(u *client.User) {
	saveName(int(u.Id), u.FirstName, u.LastName, u.Username)
}