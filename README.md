# HOW TO BUILD

go build -tags netgo -gcflags "-trimpath /go/src" -ldflags "-s -w -extldflags '-static'" github.com/paradoxgery/batch-rename
