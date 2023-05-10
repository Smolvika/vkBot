package vkontakte

import (
	"botVk/internal/currency"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

const (
	usd = "$"
	eur = "€"
	cny = "¥"
	try = "₤"
)

var nameСurrency = map[string]struct{}{"USD": {}, "EUR": {}, "CNY": {}, "TRY": {}}

var infoСurrency = map[string]struct{}{
	"Цена $ в ₽":               {},
	"Цена € в ₽":               {},
	"Цена ¥ в ₽":               {},
	"Цена ₤ в ₽":               {},
	"Изменение цены $ в ₽ и %": {},
	"Изменение цены € в ₽ и %": {},
	"Изменение цены ¥ в ₽ и %": {},
	"Изменение цены ₤ в ₽ и %": {},
}
var response []object.UsersUser

func (b *Bot) welcomeMessage(obj events.MessageNewObject, p *params.MessagesSendBuilder) {
	parameters := api.Params{
		"user_ids": obj.Message.PeerID,
	}
	err := b.vk.RequestUnmarshal("users.get", &response, parameters)
	if err != nil {
		errorsUnmarshal(err)
	}
	for _, u := range response {
		p.Message(fmt.Sprintf(`✨ %v, Вас приветствует "Бот Курс валюты" ! ✨
                  ➖Для начала работы выберите валюту, о которой вы хотите получить информацию по данным ММВБ;
                  ➖После выбора валюты Вам будут доступны кнопки с информцией, которую можно получить о валюте.
                           P.S: Вся информация берется с сайта https://bankiros.ru/currency`, u.FirstName))
	}
	_, err = b.vk.MessagesSend(p.Params)
	if err != nil {
		errorsMessage(placeSendMessageGreeting, err)
	}
}
func (b *Bot) currencyKeyboardMessage(p *params.MessagesSendBuilder) {
	p.Message("Выберите валюту:")
	p.Keyboard(getFirstLayerKeyboard())
}

func switchingСurrency(b *params.MessagesSendBuilder, str string) {
	b.Message("Выберите инфорацию,которую вы хотите получить:")
	b.Keyboard(getSecondLayerKeyboard(str))
}

func switchingInformation(b *params.MessagesSendBuilder, cur *currency.InfoCurrency, str string, errInfoCurrencyPars *error) {
	var msg string
	if *errInfoCurrencyPars == nil {
		switch str {
		case "Цена $ в ₽":
			msg = cost(usd, cur.CostUSD)
		case "Изменение цены $ в ₽ и %":
			msg = changeCost(cur.IsDecreaseUSD, usd, cur.ChangeCostRubUSD, cur.ChangeCostPrUSD)
		case "Цена € в ₽":
			msg = cost(eur, cur.CostEUR)
		case "Изменение цены € в ₽ и %":
			msg = changeCost(cur.IsDecreaseEUR, eur, cur.ChangeCostRubEUR, cur.ChangeCostPrEUR)
		case "Цена ¥ в ₽":
			msg = cost(cny, cur.CostCNY)
		case "Изменение цены ¥ в ₽ и %":
			msg = changeCost(cur.IsDecreaseCNY, cny, cur.ChangeCostRubCNY, cur.ChangeCostPrCNY)
		case "Цена ₤ в ₽":
			msg = cost(try, cur.CostTRY)
		case "Изменение цены ₤ в ₽ и %":
			msg = changeCost(cur.IsDecreaseTRY, try, cur.ChangeCostRubTRY, cur.ChangeCostPrTRY)
		}
	} else {
		msg = fmt.Sprint("Возникли проблемы с доступом к сайту биржи,вы сможете ознакомиться с курсом валюты, как только проблема будет решена. \nПриносим извинения за доставленные неудобства.")
	}
	b.Message(msg)
}

func cost(cur string, cost float64) string {
	return fmt.Sprintf("На данный момент цена 1%s = %v₽", cur, cost)
}

func changeCost(decrease bool, cur string, rub, pr float64) string {
	if decrease {
		return fmt.Sprintf("За последние 24 часа цена %s снизилась на %v₽ (%v%%)", cur, rub, pr)
	}
	return fmt.Sprintf("За последние 24 часа цена %s повысилась на %v₽ (%v%%)", cur, rub, pr)
}

func checkCurrency(cur string) bool {
	_, ok := nameСurrency[cur]
	return ok
}
func checkInfo(cur string) bool {
	_, ok := infoСurrency[cur]
	return ok
}
