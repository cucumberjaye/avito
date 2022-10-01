package apilayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const apikey = "Token"

func Convert(currency string, sum float64) (float64, error) {
	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%d", currency, int(sum))

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", apikey)
	if err != nil {
		log.Printf(err.Error())
		return sum, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf(err.Error())
		return sum, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf(err.Error())
		return sum, err
	}
	convertSum, err := layerUnmarshal(body)
	if err != nil {
		return sum, err
	}

	return convertSum, nil
}

func layerUnmarshal(body []byte) (float64, error) {
	var result ResultData
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	if result.ErrorData.Code != "" || result.Success == false {
		return 0, errors.New(result.ErrorData.Message)
	}

	return result.Result, nil
}
