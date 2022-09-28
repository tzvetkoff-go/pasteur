package hasher

// Config ...
type Config struct {
	Salt      string `yaml:"salt"`
	Alphabet  string `yaml:"alphabet"`
	MinLength int    `yaml:"min-length"`
}
