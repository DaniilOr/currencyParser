package parser

import (
	"github.com/DaniilOr/currencyParser/cmd/app/dtos"
	"encoding/json"
	"io/ioutil"
	"net/http"
)
type Parser struct {
	Url string
}

func InitService(url string) *Parser{
	return &Parser{Url: url}
}

func (s*Parser) GetUpdate() ([]*dtos.Currency, error){
	response, err := http.Get(s.Url)
	if err != nil{
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return nil, err
	}
	var Currencies []*dtos.Currency
	err = json.Unmarshal(body, &Currencies)
	if err != nil{
		return nil, err
	}
	return Currencies, nil
}