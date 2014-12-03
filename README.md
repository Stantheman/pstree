A plain pstree clone made hard to maintain. If you'd like to see the nice version, go here:

https://github.com/Stantheman/pstree/tree/c9cf104425cd1114f4ad3eb239005e655df2f3d5

I've made some changes to it since then to help speed it up. It's become less idiomatic, and more tiring to read. Along the way, I've picked up some fun benchmark numbers -- you can see the increases over many small modifications:

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

Most changes were small, and all were guided by the help of the pprof tool, including:

 * Replacing the regex for parsing /proc/pid/stat with calls to plain old string functions like index, lastindex, and fields
 * Replacing the plain old string functions with only calls to IndexByte and some slice hackery
 * Switch to using a plain old buffer instead of ioutil's ReadAll
 * Switch from using a plain old buffer to using read on the filehandle directly
 * Don't defer the filehandle close if you know you're closing it a few lines later
 * Instead of globbing on /proc/[0-9], using readdir on a proc filehandle
 * Instead of using strconv.Atoi, writing a small isInt helper method to only look at PID directories

I think I've exhausted most of the practical changes I can make since the profiler is mostly pointing to byte-to-slice conversion and file operations now.
