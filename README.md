A plain pstree clone made hard to maintain. If you'd like to see the nice version, go here:

https://github.com/Stantheman/pstree/tree/c9cf104425cd1114f4ad3eb239005e655df2f3d5

I've made some changes to it since then to help speed it up. It's become less idiomatic, and more tiring to read. Along the way, I've picked up some fun benchmark numbers:

```
BenchmarkPopulate            500           5910036 ns/op
BenchmarkPopulate            500           4122985 ns/op
BenchmarkPopulate           1000           2699748 ns/op
BenchmarkPopulate           1000           2447752 ns/op
BenchmarkPopulate           1000           2264165 ns/op
BenchmarkPopulate           1000           2061597 ns/op
BenchmarkPopulate           1000           1761342 ns/op
BenchmarkPopulate           1000           1560530 ns/op
```

Now, who needs a really fast version of a program that doesn't do everything pstree does? I'm not sure. But it's been fun trying to squeeze some speed out.

