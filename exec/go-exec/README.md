# go-exec

Build:
```
go build
```

Does not exit and hang.
Calling exit syscall directly seems not to play well with the Go runtime.
The hello world example in the same file did work.

Run:
```
./go-exec
```
