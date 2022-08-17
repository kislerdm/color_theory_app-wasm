# Color Theory App - WASM

Carnation of the [Color Theory App](https://github.com/kislerdm/color_theory_app) with two main differences:

- Logic is executed entirely on the client;
- Interface is implemented using vanilla JS instead of ReactJS.

## Assets size minification

Total volume of assets data transferred to the client over the network is an important subject for optimisation.

Note that the volume is assesses by using `wc -c`.

### .wasm file

|           Compiler            | Volume [bytes] | Comment                                                                                                        |
|:-----------------------------:|---------------:|:---------------------------------------------------------------------------------------------------------------|
|       Default Go build        |        2229903 | -                                                                                                              |
| [tinygo](https://tinygo.org/) |         516864 | -                                                                                                              |
| [tinygo](https://tinygo.org/) |         343122 | Removed dependency on `fmt`                                                                                    |
| [tinygo](https://tinygo.org/) |         157227 | Removed dependency on `fmt`<br>Build flags: `-gc=leaking -opt=2 -no-debug -panic=trap`                         |

One can see that the logic refactoring, use of different compiler with configuration tweaks leads to the binary's size reduction by the factor of **~14** (!). 
More adjustments and configuration tweaks may lead to further binary size cut. 
Although any further effort investment shall be considered carefully since its ROI may not be as high. 

In case binary size is of critical importance, `rust` could be considered as the language which may potentially yield the wasm binary of under 50kB.      

Note that the tinygo compiler does not [support](https://tinygo.org/docs/reference/lang-support/) reflection, hence the logic had to be adjusted to avoid using `encoding/json` and `encoding/csv`: 
- [Code generator for model definition](./internal/colortype/train/main.go) is used to convert the JSON model definition to the native Go struct;
- [Code generator for colors names](./internal/colorname/data/main.go) is used to convert the CSV color names map to the native Go struct.

### wasm_exec.js

| Description | Volume [bytes] | Comment |
|:-----------:|---------------:|:--------|
|   Default   |          18669 | -       |
|   TinyGo    |          16001 | -       |


## Performance optimisation

### [colorname](./internal/colorname)

```bash
go test -bench=. -benchmem ./internal/colorname 
goos: darwin
goarch: arm64
pkg: github.com/kislerdm/color_theory_app-wasm/internal/colorname
BenchmarkFindColorNameByRGB-10              6038            196192 ns/op            8046 B/op          9 allocs/op
BenchmarkFindColorNameByRGBv2-10          135073              8787 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/kislerdm/color_theory_app-wasm/internal/colorname    3.554s
```

The logic rework led to further reduction of the wasm binary size to 147791 bytes as an additional beneficial effect on top of computational and memory allocation performance improvements.

Further adjustment of the "fast" `sqrt` function by tweaking data types does not lead to performance improvement.

**Definitions**:

```go
package main

import "math"

func sqrt(v float64) float64 {
	// from quake3 inverse sqrt algorithm
	// ref: https://medium.com/@adrien.za/fast-inverse-square-root-in-go-and-javascript-for-fun-6b891e74e5a8
	const magic64 = 0x5FE6EB50C7B537A9

	n2, th := v*0.5, float64(1.5)
	b := math.Float64bits(v)
	b = magic64 - (b >> 1)
	f := math.Float64frombits(b)
	f *= th - (n2 * f * f)
	return f
}

func sqrtV2(v uint32) float32 {
	// from quake3 inverse sqrt algorithm
	// ref: https://medium.com/@adrien.za/fast-inverse-square-root-in-go-and-javascript-for-fun-6b891e74e5a8
	const magic = 0x5F375A86

	n2, th := float32(v)/2, float32(1.5)
	b := magic - (v >> 1)
	f := math.Float32frombits(b)
	f *= th - (n2 * f * f)
	return f
}
```

**Tests**:

```go
package main

import (
	"testing"
)

func benchmarkSQRT(i float64, b *testing.B) {
	for n := 0; n < b.N; n++ {
		sqrt(i)
	}
}

func benchmarkSQRTv2(i uint32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		sqrtV2(i)
	}
}

func BenchmarkSQRT1(b *testing.B)       { benchmarkSQRT(1, b) }
func BenchmarkSQRT10(b *testing.B)      { benchmarkSQRT(10, b) }
func BenchmarkSQRT100(b *testing.B)     { benchmarkSQRT(100, b) }
func BenchmarkSQRT1000(b *testing.B)    { benchmarkSQRT(1000, b) }
func BenchmarkSQRT10000(b *testing.B)   { benchmarkSQRT(10000, b) }
func BenchmarkSQRT100000(b *testing.B)  { benchmarkSQRT(100000, b) }
func BenchmarkSQRT1000000(b *testing.B) { benchmarkSQRT(1000000, b) }

func BenchmarkSQRTv2_1(b *testing.B)       { benchmarkSQRTv2(1, b) }
func BenchmarkSQRTv2_10(b *testing.B)      { benchmarkSQRTv2(10, b) }
func BenchmarkSQRTv2_100(b *testing.B)     { benchmarkSQRTv2(100, b) }
func BenchmarkSQRTv2_1000(b *testing.B)    { benchmarkSQRTv2(1000, b) }
func BenchmarkSQRTv2_10000(b *testing.B)   { benchmarkSQRTv2(10000, b) }
func BenchmarkSQRTv2_100000(b *testing.B)  { benchmarkSQRTv2(100000, b) }
func BenchmarkSQRTv2_1000000(b *testing.B) { benchmarkSQRTv2(1000000, b) }
```

Benchmark results:

```bash
go test -bench=. -benchmem .
goos: darwin
goarch: arm64
pkg: srt
BenchmarkSQRT1-10               1000000000               0.3133 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT10-10              1000000000               0.3111 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT100-10             1000000000               0.3105 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT1000-10            1000000000               0.3109 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT10000-10           1000000000               0.3114 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT100000-10          1000000000               0.3104 ns/op          0 B/op          0 allocs/op
BenchmarkSQRT1000000-10         1000000000               0.3126 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_1-10            1000000000               0.3109 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_10-10           1000000000               0.3105 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_100-10          1000000000               0.3105 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_1000-10         1000000000               0.3149 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_10000-10        1000000000               0.3114 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_100000-10       1000000000               0.3109 ns/op          0 B/op          0 allocs/op
BenchmarkSQRTv2_1000000-10      1000000000               0.3111 ns/op          0 B/op          0 allocs/op
PASS
ok      srt     5.129s
```
