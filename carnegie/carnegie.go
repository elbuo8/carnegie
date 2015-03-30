package carnegie

import (
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"time"
)

type Carnegie struct {
	Cache         *Cache
	Server        *http.Server
	CacheInterval time.Duration
}

func New(config *viper.Viper) (*Carnegie, error) {
	cache, err := NewCache(config)
	if err != nil {
		return nil, err
	}
	config.SetDefault("interval", 60*time.Second)
	carnegie := &Carnegie{
		Cache:         cache,
		CacheInterval: config.GetDuration("interval"),
	}
	config.SetDefault("address", ":8181")
	srv := &http.Server{
		Addr:    config.GetString("address"),
		Handler: http.HandlerFunc(carnegie.Handler),
	}
	srv.SetKeepAlivesEnabled(false)
	carnegie.Server = srv
	return carnegie, nil
}

func (c *Carnegie) Start() error {
	go c.UpdateCacheLoop()
	return c.Server.ListenAndServe()
}

func (c *Carnegie) RoundTrip(r *http.Request) (*http.Response, error) {
	r.RequestURI = ""
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	statusCode := res.StatusCode
	if statusCode >= 500 {
		c.Cache.Invalidate(r.Host)
	}
	return res, nil
}

func (c *Carnegie) UpdateCacheLoop() {
	ticker := time.NewTicker(c.CacheInterval)
	for {
		select {
		case <-ticker.C:
			c.Cache.LocalInventory.Purge()
		}
	}
}

func (c *Carnegie) Handler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	urls, err := c.Cache.GetAddresses(host)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(urls[0])
	proxy.Transport = c
	proxy.ServeHTTP(w, r)
	c.Cache.RotateAddresses(host)
}
