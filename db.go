package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func getLatestName(userID int) string {
	records := getNames(userID)
	if len(records) == 0 {
		return ""
	}
	return records[0].Name
}

func getNames(userID int) (n []Name) {
	records := db.ZRevRangeByScoreWithScores(getUserKey(userID), &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Val()
	for _, r := range records {
		n = append(n, Name{
			Name:     r.Member.(string),
			LastSeen: time.Unix(int64(r.Score), 0),
		})
	}
	return
}

func storeName(userID int, name string) {
	db.ZAdd(getUserKey(userID), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: name,
	})
}

func getSetLastChanged(userID int, t time.Time) *time.Time {
	i, _ := db.GetSet(getLastChangedKey(userID), t.Unix()).Int64()
	if i == 0 {
		return nil
	}
	oldT := time.Unix(i, 0)
	return &oldT
}

func storeLastChanged(userID int, t time.Time) {
	db.Set(getLastChangedKey(userID), t.Unix(), 0)
}

func getUserKey(userID int) string {
	return fmt.Sprintf("user.%d", userID)
}

func getLastChangedKey(userID int) string {
	return fmt.Sprintf("lastChanged.%d", userID)
}

type Name struct {
	Name     string
	LastSeen time.Time
}
