package api

import (
	"crypto/rand"
	"math/big"
)

type API struct {
	enabled bool
	key     string
}

func New() (api *API) {
	api = &API{
		enabled: true,
		key:     generateRandomAPI(),
	}
	return api
}

func (api *API) IsEnabled() bool {
	return api.enabled
}

func (api *API) Get() string {
	return api.key
}

func (api *API) Update() {
	api.key = generateRandomAPI()
}

func generateRandomAPI() string {
	n := 20
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		ret[i] = chars[num.Int64()]
	}
	return string(ret)
}
