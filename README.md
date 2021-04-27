# Ad Beacon URL Compression Benchmarking
This work is to investigate the performance of compression algorithms on beacon URLs.
Compressing long URLs can reduce network latency and save significant amount of memory for caching services like Redis.
The Ad beacon URLs consist of random non-repeating strings which are different to the data used in most of the compression algorithm benchmarks publicly available.
This work is done to find an algorithm that strikes a good balance in compression ratio and CPU usage on beacon URLs.
Algorithms investigated include Gzip and Snappy.

Note: In non-batch mode the data contained in the files is compressed **line by line** while in batch mode the entire data is compressed in whole as a single long string. Batch mode usually gives much better compression ratio.

## Files
* main.go: Compress data in a file using gzip or snappy. Return the total size of the compressed data.
* compress_test.go: Runs benchmark for gzip and snappy using data in file "data/urls.txt".
* data/urls.txt: Contains 70 Ad beacon URLs.
* data/long_url.txt. Contains a long Ad beacon URL whose length is 4KB.

## Usage

### Compression Ratio
To get the compression ratio of Snappy on the lines contained in a file
```go run main.go -m snappy -i ./data/urls.txt```

To get the compression ratio of Snappy on the a file
```go run main.go -m snappy -b -i ./data/urls.txt```

To get the compression ratio of Gzip on the lines contained in a file
```go run main.go -m gzip -i ./data/urls.txt```

To get the compression ratio of Gzip (level 9) on the lines contained in a file
```go run main.go -m gzip -l 9 -i ./data/urls.txt```

Gzip compression levels:
```
const (
    NoCompression      = 0
    BestSpeed          = 1
    BestCompression    = 9
    DefaultCompression = -1

    // HuffmanOnly disables Lempel-Ziv match searching and only performs Huffman
    // entropy encoding. This mode is useful in compressing data that has
    // already been compressed with an LZ style algorithm (e.g. Snappy or LZ4)
    // that lacks an entropy encoder. Compression gains are achieved when
    // certain bytes in the input stream occur more frequently than others.
    //
    // Note that HuffmanOnly produces a compressed output that is
    // RFC 1951 compliant. That is, any valid DEFLATE decompressor will
    // continue to be able to decompress this output.
    HuffmanOnly = -2
)
```

### Benchmark
To run benchmarks for Snappy and Gzip on data in "data/urls.txt"
```go test -bench=.```

## Results

### Snappy compression ratio
Snappy hardly compresses the provided URLs in non-batch mode with a compression ratio of 0.96.
It achieves a much better ratio of 0.35 in batch mode.

```
$ go run main.go -m snappy -i ./data/urls.txt
Result: totalIn=28420, totalOut=27272, ratio=0.96

$ $ go run main.go -m snappy -b -i ./data/urls.txt
Result: totalIn=28490, totalOut=9863, ratio=0.35
```

### Gzip compression ratio
The compression ratio of Gzip (default compression level) is 0.77 in non-batch mode and 0.24 in batch mode.

```
$ go run main.go -m gzip -i ./data/urls.txt
Result: totalIn=28420, totalOut=21786 ratio=0.77

$ go run main.go -m gzip -b -i ./data/urls.txt
Result: totalIn=28490, totalOut=6743, ratio=0.24
```

### Speed
Snappy performed 300 times faster than Gzip in non-batch mode on the data in file "data/urls.txt".
Snappy performed 20 times faster than Gzip in batch mode.

```
$ go test -bench=.
BenchmarkSnappy-12                 31876             36705 ns/op
BenchmarkSnappyBatch-12            59689             19848 ns/op
BenchmarkGzip-12                     126           9571068 ns/op
BenchmarkGzipBatch-12               3021            400934 ns/op
```
