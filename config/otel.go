package config

type Otel struct {
	Exporter OtelExporter `mapstructure:"exporter"`
}

type OtelExporter struct {
	Jaeger *Jaeger `mapstructure:"jaeger"`
}

type Jaeger struct {
	Endpoint string `mapstructure:"endpoint"`
}
