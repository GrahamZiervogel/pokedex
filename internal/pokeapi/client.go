package pokeapi

import (
	"net/http"
	"time"

	"github.com/GrahamZiervogel/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(httpClientTimeout, cacheReapInterval time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: httpClientTimeout,
		},
		cache: pokecache.NewCache(cacheReapInterval),
	}
}
