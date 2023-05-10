package main

import (
	"botVk/internal/vkontakte"
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("errors loading env variables: %s", err.Error())
	}
	vk := api.NewVK(os.Getenv("TOKEN"))

	group, err := vk.GroupsGetByID(api.Params{})
	if err != nil {
		log.Fatal(err)
	}

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	bot := vkontakte.NewBot(vk, lp)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := bot.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
	time.Sleep(time.Millisecond)

	lp.Shutdown()
}
