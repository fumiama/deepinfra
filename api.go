package deepinfra

import (
	"net/http"
)

const (
	APIDeepInfra = "https://api.deepinfra.com/v1/openai/chat/completions"
)

type API struct {
	api string // api to call
	key string // key in Authorization: Bearer
}

func NewAPI(api, key string) API {
	return API{api: api, key: key}
}

func (api *API) Request(model Model) (string, error) {
	req, err := http.NewRequest("POST", api.api, model.Body())
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	err = model.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	return model.Output(), nil
}
