#!/bin/bash

x86_64-w64-mingw32-dlltool -d ../client/lib/MT/libcrypto.def -l ../client/lib/MT/libcrypto.a -k
x86_64-w64-mingw32-dlltool -d ../client/lib/MT/libssl.def -l ../client/lib/MT/libssl.a -k


x86_64-w64-mingw32-g++ -o blah.exe ../client/client.cpp ../client/src/Exec.cpp ../client/src/Serialization.cpp -I../client/include -I/usr/x86_64-w64-mingw32/include -L../client/lib/MT -L/usr/x86_64-w64-mingw32/lib -lcrypto -lssl -lws2_32 -lcrypt32 -DJSON_DIAGNOSTICS=1 -static -static-libgcc -static-libstdc++ -Wl,-subsystem,console -Wl,-entry,mainCRTStartup



