package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	USD = "USDRUB_TOM"
	EUR = "EURRUB_TOM"
	CNY = "CNYRUB_TOM"
	TRY = "TRYRUB_TOM"
)

type data struct {
	cost          float64
	changeCostPr  float64
	changeCostRub float64
	isDecrease    bool
}

type infoUSD struct {
	CostUSD          float64
	ChangeCostPrUSD  float64
	ChangeCostRubUSD float64
	IsDecreaseUSD    bool
}

type infoEUR struct {
	CostEUR          float64
	ChangeCostPrEUR  float64
	ChangeCostRubEUR float64
	IsDecreaseEUR    bool
}
type infoCNY struct {
	CostCNY          float64
	ChangeCostPrCNY  float64
	ChangeCostRubCNY float64
	IsDecreaseCNY    bool
}
type infoTRY struct {
	CostTRY          float64
	ChangeCostPrTRY  float64
	ChangeCostRubTRY float64
	IsDecreaseTRY    bool
}

type InfoCurrency struct {
	infoUSD
	infoEUR
	infoCNY
	infoTRY
}

func ParsAllInfoCurrency() (InfoCurrency, error) {
	client := http.Client{Timeout: 3 * time.Second}
	usd := infoUSD{}
	err := dataCurrency(&client, USD, &usd)
	if err != nil {
		return InfoCurrency{}, err
	}
	eur := infoEUR{}
	err = dataCurrency(&client, EUR, &eur)
	if err != nil {
		return InfoCurrency{}, err
	}
	cny := infoCNY{}
	err = dataCurrency(&client, CNY, &cny)
	if err != nil {
		return InfoCurrency{}, err
	}

	try := infoTRY{}
	err = dataCurrency(&client, TRY, &try)
	if err != nil {
		return InfoCurrency{}, err
	}

	return InfoCurrency{
		infoUSD: usd,
		infoEUR: eur,
		infoCNY: cny,
		infoTRY: try,
	}, nil
}

func dataCurrency(client *http.Client, str string, cur interface{}) error {
	info, err := parsCurrency(client, str)
	if err != nil {
		return err
	}
	switch cur.(type) {
	case *infoUSD:
		cur := cur.(*infoUSD)
		*cur = infoUSD{
			CostUSD:          info.cost,
			ChangeCostPrUSD:  info.changeCostPr,
			ChangeCostRubUSD: info.changeCostRub,
			IsDecreaseUSD:    info.isDecrease,
		}
	case *infoEUR:
		cur := cur.(*infoEUR)
		*cur = infoEUR{
			CostEUR:          info.cost,
			ChangeCostPrEUR:  info.changeCostPr,
			ChangeCostRubEUR: info.changeCostRub,
			IsDecreaseEUR:    info.isDecrease,
		}
	case *infoCNY:
		cur := cur.(*infoCNY)
		*cur = infoCNY{
			CostCNY:          info.cost,
			ChangeCostPrCNY:  info.changeCostPr,
			ChangeCostRubCNY: info.changeCostRub,
			IsDecreaseCNY:    info.isDecrease,
		}
	case *infoTRY:
		cur := cur.(*infoTRY)
		*cur = infoTRY{
			CostTRY:          info.cost,
			ChangeCostPrTRY:  info.changeCostPr,
			ChangeCostRubTRY: info.changeCostRub,
			IsDecreaseTRY:    info.isDecrease,
		}

	}
	return nil
}

func parsCurrency(client *http.Client, tom string) (data, error) {
	info := data{}
	url := fmt.Sprintf("https://bankiros.ru/ajax/moex-rate-current?currency_code=%s", tom)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return info, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	res, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return info, err
	}
	jsCur1 := jsCurrency{}
	jsonErr := json.Unmarshal(body, &jsCur1)
	if jsonErr != nil {
		return info, err
	}
	for _, j := range jsCur1.Data {
		info.cost, err = strconv.ParseFloat(j.Last, 64)
		if err != nil {
			return info, err
		}
		info.changeCostRub, err = strconv.ParseFloat(strings.ReplaceAll(j.Change, "-", ""), 64)
		if err != nil {
			return info, err
		}
		info.changeCostPr, err = strconv.ParseFloat(j.ChangePercent, 64)
		if err != nil {
			return info, err
		}
		if j.Change[0] == '-' {
			info.isDecrease = true
		}
	}
	return info, nil
}
