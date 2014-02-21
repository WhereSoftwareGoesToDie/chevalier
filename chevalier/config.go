package main

type Config struct {
	Chevalier struct {
		// Just the hostname, not the port.
		ElasticsearchHost string
		// ZMQ URI to listen on.
		ListenAddress string
		// Vaultaire full read endpoint
		VaultaireReadEndpoint string
		// Vaultaire update endpoint
		VaultaireUpdateEndpoint string
	}
}
