package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"log"
	"os"
	"testing"
	
	"github.com/golang/snappy"
)

var (
	dat []string
)

func BenchmarkSnappy(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, line := range dat { 
			src := []byte(line)
			snappy.Encode(nil, src)
		}
	}
}

func BenchmarkGzip(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, line := range dat { 
			var buf bytes.Buffer
    		zw := gzip.NewWriter(&buf)
    		_, err := zw.Write([]byte(line))
			if err != nil {
				log.Fatal(err)
			}
			if err := zw.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func setup() {
	fn := "./data/urls.txt"
	file, err := os.Open(fn)
    if err != nil {
        log.Fatalf("failed to open file %s", *in)
    }
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
		dat = append(dat, scanner.Text())
    }
}

func shutdown() {
}

func TestMain(m *testing.M) {
    setup()
    code := m.Run() 
    shutdown()
    os.Exit(code)
}