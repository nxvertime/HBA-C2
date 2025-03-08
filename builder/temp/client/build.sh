#!/bin/bash

x86_64-w64-mingw32-dlltool -d lib/MT/libcrypto.def -l lib/MT/libcrypto.a -k
x86_64-w64-mingw32-dlltool -d lib/MT/libssl.def -l lib/MT/libssl.a -k



x86_64-w64-mingw32-g++ -o test.exe client.cpp \        
    -Iinclude -I/usr/x86_64-w64-mingw32/include \
    -Llib/MT -L/usr/x86_64-w64-mingw32/lib \
    -lcrypto -lssl \
    -static -static-libgcc -static-libstdc++ \
    -Wl,-subsystem,console -Wl,-entry,mainCRTStartup

