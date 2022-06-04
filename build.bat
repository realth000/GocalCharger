@echo off
echo build cmd/server...
go build --buildmode=exe -o GocalChargerServer.exe gocalcharger/cmd/server || goto fail
echo build cmd/client...
go build --buildmode=exe -o GocalChargerClient.exe gocalcharer/cmd/client || goto fail
echo build gui...
go build --buildmode=exe -o GpcalChargerGui.exe gocalcharger/gui || goto fail
echo done

pause
exit /b 0


:fail
echo build failed
pause
