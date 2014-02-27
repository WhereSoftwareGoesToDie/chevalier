package main

type Config struct {
	Chevalier struct {
		// ZMQ URI to listen on.
		ListenAddress string
	}
	Elasticsearch struct {
		// Just the hostname, not the port.
		Host     string
		Index    string
		DataType string
	}
	Vaultaire struct {
		// Vaultaire full read endpoint
		ReadEndpoint string
		// Vaultaire update endpoint
		UpdateEndpoint string
	}
}
