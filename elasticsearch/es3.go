package elasticsearch

import (
	"time"

	"github.com/alauda/bergamot/diagnose"
	elastic3 "gopkg.in/olivere/elastic.v3"
)

// ElasticConfig configuration for Elastic Search instance
type ElasticConfig struct {
	Endpoint           string
	Username           string
	Password           string
	Retries            int
	HealthCheckTimeout time.Duration
}

// GetClientOption returns a elasticSearch client options
func (c ElasticConfig) GetClientOption() []elastic3.ClientOptionFunc {
	options := make([]elastic3.ClientOptionFunc, 4, 5)
	options[0] = elastic3.SetURL(c.Endpoint)
	options[1] = elastic3.SetMaxRetries(c.Retries)
	options[2] = elastic3.SetSniff(false)
	options[3] = elastic3.SetHealthcheckTimeoutStartup(c.HealthCheckTimeout)
	if c.Username != "" {
		options = append(options, elastic3.SetBasicAuth(c.Username, c.Password))
	}
	return options
}

// ElasticSearch3Client ElasticSearch client
type ElasticSearch3Client struct {
	Client *elastic3.Client
	config ElasticConfig
}

// NewElasticSearch3Client new ES client configuration
func NewElasticSearch3Client(config ElasticConfig) (*ElasticSearch3Client, error) {
	client, err := elastic3.NewClient(
		config.GetClientOption()...,
	)

	return &ElasticSearch3Client{
		Client: client,
		config: config,
	}, err
}

// Diagnose runs a diagnose over ES connection
func (es *ElasticSearch3Client) Diagnose() diagnose.ComponentReport {
	return diagnose.SimpleDiagnose("elastic_search", func() error {
		_, _, err := es.Client.Ping(es.config.Endpoint).Do()
		return err
	})
}
