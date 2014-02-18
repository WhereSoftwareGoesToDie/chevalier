package main

type Config struct {
	Chevalier struct {
		// Just the hostname, not the port.
		ElasticsearchHost string
		// ZMQ URI to listen on.
		ListenAddress string
	}
}

