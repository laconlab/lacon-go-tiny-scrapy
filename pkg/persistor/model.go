package persistor

type Data interface {
	GetRawWebsiteAsJSON() []byte
}

type PersistorConfig struct {
	Config struct {
		Path       string `yaml:"savePath"`
		BufferSize int    `yaml:"bufferSize"`
	} `yaml:"persistor"`
}

func (s *PersistorConfig) getPath() string {
	return s.Config.Path
}

func (s *PersistorConfig) getBufferSize() int {
	return s.Config.BufferSize
}
