#include "../include/Commands.h"
#include <Windows.h>
#include <iostream>
int ExecCommandFromList(std::vector<std::string> Cmd ) {

    std::string command = "";
    for (auto& arg : Cmd) {
        command += arg;
        command += " ";
    }


    HANDLE hRead, hWrite;
    SECURITY_ATTRIBUTES sa = { sizeof(SECURITY_ATTRIBUTES), nullptr, TRUE };

    if (!CreatePipe(&hRead, &hWrite, &sa, 0)) {
        throw std::runtime_error("Error while creating pipe");
    }

    // Configure la redirection pour le processus enfant
    STARTUPINFO si = { sizeof(STARTUPINFO) };
    si.dwFlags = STARTF_USESTDHANDLES;
    si.hStdOutput = hWrite;
    si.hStdError = hWrite;

    // Informations sur le processus
    PROCESS_INFORMATION pi = { 0 };
    
    std::string cmd = "cmd.exe /C " + command; // Utilisation de cmd.exe pour exécuter la commande

    
    if (!CreateProcessA(nullptr, cmd.data(), nullptr, nullptr, TRUE, CREATE_NO_WINDOW, nullptr, nullptr, &si, &pi)) {
        CloseHandle(hRead);
        CloseHandle(hWrite);
        throw std::runtime_error("Error while creating the process");
    }

    
    CloseHandle(hWrite);

    
    char buffer[128];
    DWORD bytesRead;
    std::string result;
    while (ReadFile(hRead, buffer, sizeof(buffer) - 1, &bytesRead, nullptr) && bytesRead > 0) {
        buffer[bytesRead] = '\0'; 
        result += buffer;
    }

    CloseHandle(hRead);
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);


    std::cout << "[EXECCMD] RESULT: " << result << std::endl;


}