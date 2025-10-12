package initialization

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tienhai2808/ecom_go/internal/config"
)

func InitElasticsearch(cfg *config.Config) (*elasticsearch.TypedClient, error) {
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("kết nối tới Elasticsearch thất bại: %w", err)
	}

	return es, nil
}
