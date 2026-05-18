package elasticsearch

import (
	"log"

	"github.com/PuppyNote/puppynote-server-golang/config"
	"github.com/elastic/go-elasticsearch/v8"
)

var Client *elasticsearch.Client

func Connect() {
	if config.AppConfig.Elasticsearch.Host == "" {
		log.Println("[ES] ES_HOST 미설정, Elasticsearch 비활성화")
		return
	}

	cfg := elasticsearch.Config{
		Addresses: []string{config.AppConfig.Elasticsearch.Host},
		Username:  config.AppConfig.Elasticsearch.Username,
		Password:  config.AppConfig.Elasticsearch.Password,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("[ES] 연결 실패 (비치명적): %v", err)
		return
	}
	Client = client
	log.Printf("[ES] Elasticsearch 연결: %s", config.AppConfig.Elasticsearch.Host)
}
