package metrics

import (
	"log"
	"github.com/Shopify/toxiproxy/v2/client"
)

func SetupToxiProxy(apiAddr string) (*toxiproxy.Client, error) {
	toxiClient := toxiproxy.NewClient(apiAddr)
	log.Println("ToxiProxy client created")

	// Create the proxy 
	_, err := toxiClient.CreateProxy("minio_proxy", "0.0.0.0:8666", "minio:9000")
	if err != nil {
		log.Printf("Failed to create proxy: %v", err)
		return nil, err
	}
	
	log.Println("Proxy created successfully on port 8666")
	return toxiClient, nil
}
