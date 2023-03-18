all: zapi_unix.go
#all: zapi_windows.go zapi_unix.go

#zapi_windows.go: api.go
#	$(GOPATH)/src/pkg/syscall/mksyscall_windows.pl $^ \
#		| gofmt \
#		> $@

zapi_unix.go: api.go
	./mksyscall_unix.pl $^ \
		| gofmt \
		> $@
