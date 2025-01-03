# Find function without tracing implemented

Use this regex in VSCode search:

```regex
func.*\) \{\n\t+(?!ctx)
```

Include following files:

```regex
*.go
```

Exclude following files:

```regex
*_test.go,*_seeder.go,*_factory.go,mocks/*,*.pb.go
```
