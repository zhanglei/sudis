package main

//go get -u github.com/shuLhan/go-bindata/cmd/go-bindata

//go:generate go-bindata -pkg http -o master/server/http/http_static_bindata_assets.go webui/dist/...
