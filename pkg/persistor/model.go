package persistor

type StoreConfig struct {
	Config struct {
		SavePath    string `yaml:"savePath"`
		NamePattern string `yaml:"namePattern"`
	} `yaml:"persistor"`
}

func (s *StoreConfig) getSavePath() string {
	return s.Config.SavePath
}

func (s *StoreConfig) getNamePattern() string {
	return s.Config.NamePattern
}

type LoadConfig struct {
	Config struct {
		LoadPath   string `yaml:"savePath"`
		BufferSize int    `yaml:"bufferSize"`
	} `yaml:"persistor"`
}

func (s *LoadConfig) getLoadPath() string {
	return s.Config.LoadPath
}

func (s *LoadConfig) getBufferSize() int {
	return s.Config.BufferSize
}

type Page interface {
	GetId() int
	GetName() string
	GetUrl() string
	GetContent() []byte
}

type PageImpl struct {
	Id      int    `json:"id"`
	Url     string `json:"url"`
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

func (h *PageImpl) GetUrl() string {
	return h.Url
}

func (h *PageImpl) GetId() int {
	return h.Id
}

func (h *PageImpl) GetName() string {
	return h.Name
}

func (h *PageImpl) GetContent() []byte {
	return h.Content
}
