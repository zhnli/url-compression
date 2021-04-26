# Ad Beacon URL Compression Benchmarking
Investigate the performance of compression algorithms for beacon URLs.
Beacon redirect service has increased the load and memory usage of the Redis cluster.
The Ad beacon URLs consist of random non-repeating strings which are different to the data used in most of the compression algorithm benchmarks publicly available.
This work is done to find an algorithm that has good balance in compression ratio and CPU usage on beacon URLs.
Algorithms investigated include Gzip and Snappy.

Note: The data contained in the files are compressed **line by line** which is different to compressing the entire data as a whole. It is done this way to simulate how compression is performed in our applications.

## Files
* main.go: Compress data in a file using gzip or snappy. Return the total size of the compressed data.
* compress_test.go: Runs benchmark for gzip and snappy using data in file "data/urls.txt".
* data/urls.txt: Contains 70 Ad beacon URLs.
* data/long_url.txt. Contains a long Ad beacon URL whose length is 4KB.

## Usage

### Compression Ratio
To get the compression ratio of Snappy on the lines contained in a file
```go run main.go -m snappy -i ./data/urls.txt```

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
Snappy hardly compresses the provided URLs with a compression ratio of 0.96.

```
$ go run main.go -m snappy -i ./data/urls.txt 
Result: totalIn=28420, totalOut=27272
```

### Gzip compression ratio
The compression ratio of Gzip (default compression level) is around 0.75 on mixed URLs.
The ratio is slightly higher with longer URLs. 
```
$ go run main.go -m gzip -i ./data/urls.txt 
Result: totalIn=28420, totalOut=21786
$ go run main.go -m gzip -i ./data/urls.txt 
Result: totalIn=3992, totalOut=2805
```

### Speed
Snappy performed 300 times faster than Gzip on the data in file "data/urls.txt".
```
$ go test -bench=.
BenchmarkSnappy-12         31886             36592 ns/op
BenchmarkGzip-12             127           9385610 ns/op
```
