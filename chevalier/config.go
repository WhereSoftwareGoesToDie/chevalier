package main

type Config struct {
	Chevalier struct {
		// In host:port format.
		ElasticsearchEndpoint string
		// ZMQ URI to listen on.
		ListenAddress string
	}
}

