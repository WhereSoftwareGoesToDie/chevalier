package main

type Config struct {
	Chevalier struct {
		// ZMQ URI to listen on.
		ListenAddress string
		// Log level (info, debug, et cetera)
		LogLevel string
	}
	Elasticsearch struct {
		// Just the hostname, not the port.
		Host         string
		Index        string
		DataType     string
		MaxConns     int
		RetrySeconds int
	}
	Vaultaire struct {
		// Vaultaire full read endpoint
		ReadEndpoint string
		// Vaultaire update endpoint
		UpdateEndpoint string
		Origins        []string
	}
	Indexer struct {
		// Maximum number of coroutines to use for indexing.
		// Note that the Elasticsearch writer will use MaxConns
		// threads of its own.
		Parallelism uint
	}
}
