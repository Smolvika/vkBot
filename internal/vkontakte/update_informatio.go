package vkontakte

import (
	"botVk/internal/currency"
	"context"
	"log"
	"time"
)

func updateInfoAboutBitcoin(ctx context.Context, currencyNow *currency.InfoCurrency, errInfoCurrencyNowPars *error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			*currencyNow, *errInfoCurrencyNowPars = currency.ParsAllInfoCurrency()
			if *errInfoCurrencyNowPars != nil {
				log.Println("Problem with parsing currency information: ", *errInfoCurrencyNowPars)
			}
		}
	}
}
