package stresstest

import (
	"errors"
	"net/url"
)

type Config struct {
	url         string
	requests    int
	concurrency int
}

func NewConfig(url string, requests int, concurrency int) (*Config, error) {
	config := &Config{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
	}

	err := config.validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.url == "" {
		return errors.New("URL cannot be empty")
	}

	parsedUrl, err := url.ParseRequestURI(c.url)
	if err != nil {
		return errors.New("Url must be a valid URI")
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return errors.New("Url must start with http:// or https://")
	}

	if c.requests <= 0 {
		return errors.New("Requests must be greater than 0")
	}
	if c.concurrency <= 0 {
		return errors.New("Concurrency must be greater than 0")
	}
	return nil
}
