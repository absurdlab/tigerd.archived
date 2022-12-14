package authorize

// ProviderProperties is the configuration properties for a single provider.
type ProviderProperties struct {
	// Key is the unique identifier of this provider.
	Key string `json:"key" yaml:"key"`
	// Address is the <host>:<port> address to connect to the provider. Non-localhost is supported, however, will
	// print a WARNING message to console.
	Address string `json:"address" yaml:"address"`
}
