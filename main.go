package main
import (
    "fmt"
    "flag"
    "net/http"
    "os"
    "io"
    "io/ioutil"
    "time"
)

type responseInfo struct {
    status int
    bytes int64
    duration time.Duration
}

type SummaryInfo struct {
    requested int64
    responsed int64
}

func main() {
    fmt.Println("Hello from my app")
    requests := flag.Int64("n", 1, "Number of requests to perform")
    concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
    
    fmt.Println(requests, concurrency)

    flag.Parse()

    if flag.NArg() == 0 || *requests == 0 || *requests < *concurrency {
        flag.PrintDefaults()
        os.Exit(-1)
    }

    link := flag.Arg(0)
    summary := SummaryInfo{}
    c := make(chan responseInfo)
    for i := int64(0); i < *concurrency; i ++ {
        summary.requested++
        go checkLink(link, c)
    }

    for response := range c {
        if summary.responsed < *requests {
            summary.requested++
            go checkLink(link, c)
        }
        summary.responsed++
        fmt.Println(response)
        if summary.responsed == summary.requested {
            break
        }
    }
}

func checkLink(link string, c chan responseInfo) {
    start := time.Now()
    res, err := http.Get(link)
    if err != nil {
        panic(err)
    }
    read, _ := io.Copy(ioutil.Discard, res.Body)
    c <- responseInfo {
        status: res.StatusCode,
        bytes: read,
        duration: time.Now().Sub(start),
    }
}