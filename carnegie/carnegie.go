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
	Config        *viper.Viper
	Started       bool
}

func New(config *viper.Viper) (*Carnegie, error) {
	cache, err := NewCache(config)
	if err != nil {
		return nil, err
	}
	config.SetDefault("interval", 60*time.Second)
	config.SetDefault("address", ":8181")

	carnegie := Carnegie{
		Cache:         cache,
		CacheInterval: config.GetDuration("interval"),
		Server: &http.Server{
			Addr: config.GetString("address"),
		},
		Config:  config,
		Started: false,
	}

	carnegie.Server.Handler = http.HandlerFunc(carnegie.Handler)
	carnegie.Server.SetKeepAlivesEnabled(false)

	return &carnegie, nil
}

// BILL: What happens if I call Start 1000 times.
func (c *Carnegie) Start() error {
	if c.Started {
		return nil
	}
	c.Started = true
	go c.UpdateCacheLoop()
	if certFile, keyFile := c.Config.GetString("cert"), c.Config.GetString("key"); certFile != "" && keyFile != "" {
		go c.Server.ListenAndServeTLS(certFile, keyFile)
	}
	return c.Server.ListenAndServe()
}

func (c *Carnegie) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 500 {
		c.Cache.Invalidate(r.Host)
	}

	return res, nil
}

func (c *Carnegie) UpdateCacheLoop() {
	ticker := time.NewTicker(c.CacheInterval)
	for {
		<-ticker.C
		c.Cache.LocalInventory.Purge()
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
