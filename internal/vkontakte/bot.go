package vkontakte

import (
	"botVk/internal/currency"
	"context"
	"errors"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
)

type Bot struct {
	vk *api.VK
	lp *longpoll.LongPoll
}

func NewBot(vk *api.VK, lp *longpoll.LongPoll) *Bot {
	return &Bot{
		vk: vk,
		lp: lp,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	var errInfoCurrencyNowPars error
	currencyNow, err := currency.ParsAllInfoCurrency()
	if err != nil {
		return errors.New(fmt.Sprintf("Problem with parsing currency information %v", err))
	}
	go updateInfoAboutBitcoin(ctx, &currencyNow, &errInfoCurrencyNowPars)

	b.botHandler(&currencyNow, &errInfoCurrencyNowPars)
	log.Println("Start Long Poll")
	if err := b.lp.Run(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) botHandler(cur *currency.InfoCurrency, errInfoCurrencyPars *error) {
	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
		p := params.NewMessagesSendBuilder()
		p.RandomID(0)
		p.PeerID(obj.Message.PeerID)
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		switch {
		case obj.Message.Text == "Начать" || obj.Message.Text == "/старт":
			b.welcomeMessage(obj, p)
			fallthrough
		case obj.Message.Text == "К валюте":
			b.currencyKeyboardMessage(p)
		case checkCurrency(obj.Message.Text):
			switchingСurrency(p, obj.Message.Text)
		case checkInfo(obj.Message.Text):
			switchingInformation(p, cur, obj.Message.Text, errInfoCurrencyPars)
		default:
			log.Printf("A nonexistent command was received %s", obj.Message.Text)
			p.Message("Я не понимаю эту команду.\nНачать общение с начала: /старт")
		}
		_, err := b.vk.MessagesSend(p.Params)
		if err != nil {
			errorsMessage(placeSendMessage, err)
		}
	})
}
