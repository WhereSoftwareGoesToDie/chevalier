package main

type Config struct {
	Chevalier struct {
		// ZMQ URI to listen on.
		ListenAddress string
		// ZMQ URI to receive/respond to status requests.
		StatusAddress string
		// Log level (info, debug, et cetera)
		LogLevel string
	}
	Elasticsearch struct {
		// Just the hostname, not the port.
		Host  string
		Index string
		// Index name for chevalier-related metadata
		MetadataIndex string
		DataType      string
		MaxConns      int
		RetrySeconds  int
	}
	Vaultaire struct {
		// Vaultaire contents daemon endpoint.
		ContentsEndpoint string
		Origins          []string
	}
	Indexer struct {
		// Maximum number of coroutines to use for indexing.
		// Note that the Elasticsearch writer will use MaxConns
		// threads of its own.
		Parallelism uint
	}
}
