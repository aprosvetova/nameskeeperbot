package main

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
)

func listenTdlib() {
	authorizer := client.BotAuthorizer(cfg.Token)
	authorizer.TdlibParameters <- &client.TdlibParameters{
		DatabaseDirectory:      "/data/database",
		UseChatInfoDatabase:    true,
		ApiId:                  int32(cfg.TdlibApiID),
		ApiHash:                cfg.TdlibApiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Ebik",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 0,
	})

	tdlibClient, err := client.NewClient(authorizer, logVerbosity)
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
