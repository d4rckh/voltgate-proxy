package main

import (
	"log"
	"net/http"
	"voltgate-proxy/caching"
	"voltgate-proxy/config"
	"voltgate-proxy/handler"
	"voltgate-proxy/management"
	"voltgate-proxy/monitoring"
	"voltgate-proxy/proxy"
	"voltgate-proxy/ratelimiting"
)

func main() {
	println("            .__   __                __          \n___  ______ |  |_/  |_  _________ _/  |_  ____  \n\\  \\/ /  _ \\|  |\\   __\\/ ___\\__  \\\\   __\\/ __ \\ \n \\   (  <_> )  |_|  | / /_/  > __ \\|  | \\  ___/ \n  \\_/ \\____/|____/__| \\___  (____  /__|  \\___  >\n                     /_____/     \\/          \\/ ")

	proxyServer := proxy.NewProxyServer()
	initialConfig, rateLimitRules, cacheRules := config.LoadConfig(proxyServer, "config.yaml")

	monitoring.InitMetrics()

	switch initialConfig.RateLimitAppConfig.Storage {
	case "redis":
		proxyServer.RateLimiterStorage = ratelimiting.MakeRedisRateLimiterStorage(initialConfig.Storage.Redis)
	case "memory":
		proxyServer.RateLimiterStorage = ratelimiting.MakeInMemoryRateLimiterStorage()
	}

	switch initialConfig.CacheAppConfig.Storage {
	case "redis":
		proxyServer.CacherStorage = caching.MakeRedisCacherStorage(initialConfig.Storage.Redis)
	case "memory":
		proxyServer.CacherStorage = caching.MakeInMemoryCacherStorage()
	}

	if len(rateLimitRules.EndpointRateLimitRules) != 0 && proxyServer.RateLimiterStorage == nil {
		log.Printf("ERROR: Found rate limiting rules but rate limiter storage not initialized, please define storage in config")
		panic("")
	}

	if len(cacheRules.EndpointCacheRules) != 0 && proxyServer.CacherStorage == nil {
		log.Printf("ERROR: Found cache rules but cacher storage not initialized, please define storage in config")
		panic("")
	}

	if initialConfig.ReloadConfigInterval != 0 {
		log.Println("Config reloading enabled, interval set to", initialConfig.ReloadConfigInterval, "seconds")
		go config.ReloadConfig(proxyServer, initialConfig.ReloadConfigInterval, "config.yaml")
	}

	go management.StartManagementServer(initialConfig)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handler.HandleRequest(proxyServer, &rateLimitRules, &cacheRules, writer, request)
	})

	log.Println("Proxy server started on", initialConfig.ProxyConfig.Address)
	log.Fatal(http.ListenAndServe(initialConfig.ProxyConfig.Address, nil))
}
