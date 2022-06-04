@echo off
echo build cmd/server...
go build --buildmode=exe -o GocalChargerServer.exe ./cmd/server/server.go || goto fail
echo build cmd/client...
go build --buildmode=exe -o GocalChargerClient.exe ./cmd/client/client.go || goto fail
echo build gui...
go build --buildmode=exe -o GpcalChargerGui.exe ./gui/main.go || goto fail
echo done

pause
exit /b 0


:fail
echo build failed
pause
