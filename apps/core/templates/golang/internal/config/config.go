package config

type Config struct {
	AppName string
	Port    string

	OTLPEndpoint string
}

func Load() Config {
	return Config{
		AppName:      "{{.ServiceName}}",
		Port:         "8080",
		OTLPEndpoint: "otel-collector.observability.svc:4317",
	}
}
