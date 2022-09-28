package webserver

// Config ...
type Config struct {
	StaticPath      string `yaml:"static-path"`
	TemplatesPath   string `yaml:"templates-path"`
	ListenAddress   string `yaml:"listen-address"`
	ProxyHeader     string `yaml:"proxy-header"`
	TLSCert         string `yaml:"tls-cert"`
	TLSKey          string `yaml:"tls-key"`
	RelativeURLRoot string `yaml:"relative-url-root"`
}
