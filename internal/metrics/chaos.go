package metrics

import "github.com/Shopify/toxiproxy/v2/client"

func ControlChaos(proxyClient *toxiproxy.Client, chaosType string, stream string, toxicity float64, attributes map[string]interface{}) error {

	proxyClient.RemoveToxic("active_chaos")

	options := &toxiproxy.ToxicOptions{
		Name:       "slow_minio",      // The name of this toxic
		Type:       chaosType,         // The type (latency, bandwidth, etc.)
		Stream:     stream,            // Direction
		Toxicity:   toxicity,          // 100% of traffic
		Attributes: toxiproxy.Attributes(attributes),
	}

	_, err := proxyClient.AddToxic(options)
	return err
}
