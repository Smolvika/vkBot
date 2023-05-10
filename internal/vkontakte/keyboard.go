package vkontakte

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/object"
)

var firstLayerKeyboard = object.NewMessagesKeyboard(false)
var secondUSDLayerKeyboard = object.NewMessagesKeyboard(false)
var secondEURLayerKeyboard = object.NewMessagesKeyboard(false)
var secondCNYLayerKeyboard = object.NewMessagesKeyboard(false)
var secondTRYLayerKeyboard = object.NewMessagesKeyboard(false)

func getFirstLayerKeyboard() *object.MessagesKeyboard {
	if len(firstLayerKeyboard.Buttons) == 0 {
		row1 := firstLayerKeyboard.AddRow()
		row1.AddTextButton("USD", "", "primary")
		row1.AddTextButton("EUR", "", "primary")
		row2 := firstLayerKeyboard.AddRow()
		row2.AddTextButton("CNY", "", "primary")
		row2.AddTextButton("TRY", "", "primary")
	}
	return firstLayerKeyboard
}

func getSecondLayerKeyboard(cur string) *object.MessagesKeyboard {
	switch cur {
	case "USD":
		keyboard(secondUSDLayerKeyboard, usd)
		return secondUSDLayerKeyboard
	case "EUR":
		keyboard(secondEURLayerKeyboard, eur)
		return secondEURLayerKeyboard
	case "CNY":
		keyboard(secondCNYLayerKeyboard, cny)
		return secondCNYLayerKeyboard
	case "TRY":
		keyboard(secondTRYLayerKeyboard, try)
		return secondTRYLayerKeyboard
	}
	return nil
}
func keyboard(secondLayerKeyboard *object.MessagesKeyboard, cur string) *object.MessagesKeyboard {
	if len(secondLayerKeyboard.Buttons) == 0 {
		row1 := secondLayerKeyboard.AddRow()
		row1.AddTextButton(fmt.Sprintf("Цена %s в ₽", cur), "", "primary")
		row2 := secondLayerKeyboard.AddRow()
		row2.AddTextButton(fmt.Sprintf("Изменение цены %s в ₽ и %%", cur), "", "primary")
		row3 := secondLayerKeyboard.AddRow()
		row3.AddTextButton("К валюте", "", "negative")
	}
	return secondLayerKeyboard
}
