#SHELL=C:/Windows/System32/cmd.exe
SHELL=D:/dev/msys64/usr/bin/bash.exe

watch:
	gobble watch app/build

dev:build
	npm start

build:
	@# flatc -g query_result.fbs request.fbs response.fbs
	@# flatc -s -o ./app/js/flatapi query_result.fbs request.fbs response.fbs
	go build -o ./plugins/golaunch-programs/golaunch-programs.exe ./plugins/golaunch-programs
	go build -o ./plugins/golaunch-process-killer/golaunch-process-killer.exe ./plugins/golaunch-process-killer

# rebrand:
# 	rcedit electron\electron.exe --set-icon app\icon.ico --set-version-string "CompanyName" "Outroot" --set-version-string "FileDescription" "GoLaunch application launcher" --set-version-string "LegalCopyright" "Copyright (C) 2015 Kevin Darlington. All rights reserved." --set-version-string "ProductName" "GoLaunch" --set-file-version "0.0.1" --set-product-version "0.0.1"

# bundle:
# 	electron-packager . Playback --platform=win32 --arch=ia32 --version=0.27.2 --icon=icon.ico
