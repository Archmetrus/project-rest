package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Gateway forwards HTTP and JSON unchanged to the owning service.
func Gateway(userAddress, hrAddress string) (http.Handler, error) {
	mux := http.NewServeMux()
	for path, address := range map[string]string{"/users": userAddress, "/leaves": hrAddress} {
		target, err := url.Parse("http://" + address)
		if err != nil || target.Host == "" {
			return nil, fmt.Errorf("invalid service address: %s", address)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("upstream error: %v", err)
			respond(w, http.StatusBadGateway, map[string]string{"error": "upstream service unavailable"})
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			proxy.ServeHTTP(w, r.WithContext(ctx))
		})
		mux.Handle(path, handler)
		mux.Handle(path+"/", handler)
	}
	return mux, nil
}
