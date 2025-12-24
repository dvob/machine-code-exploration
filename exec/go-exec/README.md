# go-exec

Build:
```
go build
```

Direct calling of generated code without CGO is probably a very bad idea.

Wazero does that, but you have to take into account all sorts of constraints: https://github.com/wazero/wazero/blob/df9f68620da407c1bcaa043d2c302f2d5cb43acc/site/content/docs/how_do_compiler_functions_work.md

Run:
```
./go-exec
```

I tested this with Go 1.25.3. As this relies on the internal structure of funcval, it's possible that this will break at some point in the future.
