package crawler

type downloadStatus int

const (
	ok downloadStatus = iota
	skip
	retry
	undefined
)

func getDownloadStatusFromStatusCode(status int) downloadStatus {
	if status >= 500 {
		return retry
	} else if status >= 200 && status <= 300 {
		return ok
	} else {
		return skip
	}
}
