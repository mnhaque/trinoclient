package config

import "trino-sdk/trinodriver"

type Configuration struct {
	Uri        string
	ClientName string
	Config     trinodriver.Config
}

func New(uri string, clientName string) *Configuration {
	config := &Configuration{Uri: uri, ClientName: clientName}
	config.Config = trinodriver.Config{ServerURI: config.Uri, Catalog: "hive", Schema: "default", Source: clientName}
	return config
}

func (c *Configuration) GetDSN() (string, error) {
	return c.Config.FormatDSN()
}
