package crawler

type downloadRequest interface {
    getUrl() string
}
