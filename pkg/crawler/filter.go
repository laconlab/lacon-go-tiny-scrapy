package crawler

//import "fmt"

/*

func NewFilter(reqs chan downloadRequest) chan downloadRequest {
    output := make(chan downloadRequest)

    for i := 0; i < 10; i++ {
        go filter(reqs, output)
    }

    return output
}

func filter(reqs chan downloadRequest, output chan downloadRequest) {
    for req := range reqs {
        url := req.getUrl()
        fmt.Print(url)
    }
}
*/
