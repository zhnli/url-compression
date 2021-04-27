package main

import (
    "bufio"
    "bytes"
	"compress/gzip"
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"

	"github.com/golang/snappy"
)

var (
	in = flag.String("i", "", "input filename")
    method = flag.String("m", "", "compress method <snappy | gzip>")
    level = flag.Int("l", -1, "gzip compression level (default level -1)")
    batch = flag.Bool("b", false, "batch processing mode")
)

func compressSnappy(dat []byte) int {
    //src := []byte(dat)
    encoded := snappy.Encode(nil, dat)
    return len(encoded)
}

func compressGzip(dat []byte) int{
    var buf bytes.Buffer
    zw, _ := gzip.NewWriterLevel(&buf, *level)
    
    _, err := zw.Write(dat)
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
    return len(buf.Bytes())
}

func processLines(inFile string, method string) (int, int) {
    file, err := os.Open(inFile)
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
        if method == "snappy" {
            totalOut += compressSnappy([]byte(line))
        } else {
            totalOut += compressGzip([]byte(line))
        }
    }
    return totalIn, totalOut
}

func processBatch(inFile string, method string) (int, int) {
    dat, err := ioutil.ReadFile(inFile)
    if err != nil {
        log.Fatalf("failed to open file %s", *in)
    }
    totalOut := 0
    if method == "snappy" {
        totalOut = compressSnappy(dat)
    } else {
        totalOut = compressGzip(dat)
    }
    return len(dat), totalOut
}

func main() {
    flag.Parse()

    if len(*in) == 0 || len(*method) == 0 {
        fmt.Println("Usage: <cmd> -m <snappy|gzip> -i input_file")
        panic(errors.New("Argument missing!"))
    }

    totalIn, totalOut := 0, 0
    if *batch {
        totalIn, totalOut = processBatch(*in, *method)
    } else {
        totalIn, totalOut = processLines(*in, *method)
    }

    log.Printf("Result: totalIn=%d, totalOut=%d, ratio=%.2f", totalIn, totalOut, float32(totalOut) / float32(totalIn))
}
