@echo build linux
@cd ..
@if not exist ..\forestAdminBin\server mkdir ..\forestAdminBin\server

@set GOARCH=amd64
@set GOOS=linux
@go build

move go-admin ../forestAdminBin/server
