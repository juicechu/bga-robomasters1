### Unity Bridge

The Unity Bridge wraps up a Windows DLL and, due to that, only works on Windows 
(more specifically, Windows 64 bits).

It also makes use of CGO so one needs both Go(https://www.golang.org) installed and a C compiler (Mingw-w64 is the recommended one and can be downloaded from http://mingw-w64.org/doku.php).

###### Compiling on Linux

Due to what was mentioned above cross-compiling it inside Linux requires Go and a Windows C cross-compiler. Mingw-w64 is also the recommended one and can usually be installed directly with your distribution's package manager.

After installing Go and Mingw-w64, the code can be build with:

CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build

###### Running on Linux

The result of cross-compilation is a Windows executable. To run it one can use Wine (which can also be installed through the distribution package manager) and, as far as I can tell, it seems to work as expected. Disabling Wine logging is recommended to remove unnecessary noise:

WINEDEBUG=-all wine program.exe

