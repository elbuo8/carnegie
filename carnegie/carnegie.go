package carnegie

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// Carnegie is a load balancer for dynamic VHOST inventories.
type Carnegie struct {
	Cache         *Cache
	Server        *http.Server
	CacheInterval time.Duration
	Config        *viper.Viper
	Started       bool
	Verbose       bool
	Logger        *log.Logger
}

// New returns a new Carnegie with the provided configuration.
func New(config *viper.Viper) (*Carnegie, error) {
	cache, err := NewCache(config)
	if err != nil {
		return nil, err
	}
	config.SetDefault("interval", 60*time.Second)
	config.SetDefault("address", ":8181")
	config.SetDefault("verbose", true)
	config.SetDefault("log", "stdout")

	carnegie := Carnegie{
		Cache:         cache,
		CacheInterval: config.GetDuration("interval"),
		Server: &http.Server{
			Addr: config.GetString("address"),
		},
		Config:  config,
		Started: false,
		Verbose: config.GetBool("verbose"),
	}

	switch config.GetString("log") {
	case "stdout":
		carnegie.Logger = log.New(os.Stdout, "carnegie: ", log.Lmicroseconds)
	}

	carnegie.Server.Handler = http.HandlerFunc(carnegie.handler)
	carnegie.Server.SetKeepAlivesEnabled(false)

	return &carnegie, nil
}

// Start will start the cache updating loop as well as an HTTP listener.
// If TLS information is provided, will launch a TLS listener.
func (c *Carnegie) Start() error {
	if c.Started {
		return nil
	}
	c.Started = true
	go c.updateCacheLoop()
	if certFile, keyFile := c.Config.GetString("cert"), c.Config.GetString("key"); certFile != "" && keyFile != "" {
		go c.Server.ListenAndServeTLS(certFile, keyFile)
	}
	return c.Server.ListenAndServe()
}

// RoundTrip returns the request performed to the VHOST backend.
// If an error occurred, the VHOST will be invalidated.
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

func (c *Carnegie) updateCacheLoop() {
	ticker := time.NewTicker(c.CacheInterval)
	for {
		<-ticker.C
		c.Cache.LocalInventory.Purge()
	}
}

func (c *Carnegie) handler(w http.ResponseWriter, r *http.Request) {
	if c.Verbose {
		c.Logger.Printf("%s %s %s", r.Method, r.URL, r.Host)
	}
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
