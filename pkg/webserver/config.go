package webserver

// Config ...
type Config struct {
	Assets AssetsConfig `yaml:"assets"`
	Listen ListenConfig `yaml:"listen"`
	Proxy  ProxyConfig  `yaml:"proxy"`
}

// AssetsConfig ...
type AssetsConfig struct {
	RelativeURLRoot string `yaml:"relative-url-root"`
	StaticPath      string `yaml:"static-path"`
	TemplatesPath   string `yaml:"templates-path"`
}

// ListenConfig ...
type ListenConfig struct {
	Address string `yaml:"address"`
	TLSCert string `yaml:"tls-cert"`
	TLSKey  string `yaml:"tls-key"`
}

// ProxyConfig ...
type ProxyConfig struct {
	Enabled    bool     `yaml:"enabled"`
	IPHeader   string   `yaml:"ip-header"`
	Loopback   bool     `yaml:"loopback"`
	LinkLocal  bool     `yaml:"link-local"`
	Private    bool     `yaml:"private"`
	UnixSocket bool     `yaml:"unix-socket"`
	TrustedIPs []string `yaml:"trusted-ips"`
}
