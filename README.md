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
