package main

import (
    "bufio"
    "bytes"
	"compress/gzip"
    "errors"
    "flag"
    "fmt"
    "log"
    "os"

	"github.com/golang/snappy"
)

var (
	in = flag.String("i", "", "input filename")
    method = flag.String("m", "", "compress method <snappy | gzip>")
    level = flag.Int("l", -1, "gzip compression level (default level -1)")
)

func compressSnappy(dat string) int {
    src := []byte(dat)
    encoded := snappy.Encode(nil, src)
    return len(encoded)
}

func compressGzip(dat string) int{
    var buf bytes.Buffer
    zw, _ := gzip.NewWriterLevel(&buf, *level)
    
    _, err := zw.Write([]byte(dat))
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
    return len(buf.Bytes())
}

func main() {
    flag.Parse()

    if len(*in) == 0 || len(*method) == 0 {
        fmt.Println("Usage: <cmd> -m <snappy|gzip> -i input_file")
        panic(errors.New("Argument missing!"))
    }

    file, err := os.Open(*in)
    if err != nil {
        log.Fatalf("failed to open file %s", *in)
    }
    defer file.Close()

    totalIn, totalOut := 0, 0
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
        line := scanner.Text()
        totalIn += len(line)
        if *method == "snappy" {
            totalOut += compressSnappy(line)
        } else {
            totalOut += compressGzip(line)
        }
    }

    log.Printf("totalIn=%d, totalOut=%d", totalIn, totalOut)
}
