package main

type Config struct {
	// In host:port format.
	ElasticsearchEndpoint string
	// ZMQ URI to listen on.
	ListenAddress string
}

