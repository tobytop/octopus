
build protoc plugin (in Windows)
```
go build -o protoc-gen-meun.exe .
```

build proto by using protoc (make sure the out_path value and go_out value are same)
```
protoc --meun_out=out_path=../proto,project_name=octopus:. --go_out=../proto --proto_path=. *.proto
```