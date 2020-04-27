package client

import "github.com/elastic/go-elasticsearch/v8"

var EsClient *elasticsearch.Client

func init() {
	var err error
	config := elasticsearch.Config{}
	config.Addresses = []string{"http://127.0.0.1:9200"}
	EsClient, err = elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
}
