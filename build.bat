del chatplayer.exe

git rev-parse --short HEAD > %TEMP%\git-rev.txt

set /P VER1=<version.txt
set /P VER2=<%TEMP%\git-rev.txt

go build -ldflags "-X main.programVersion=%VER1%-%VER2%"
