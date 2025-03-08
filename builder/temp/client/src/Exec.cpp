#include "../include/Commands.h"


std::string execCommand(const std::string& command) {
    std::array<char, 128> buffer;
    std::string result;

    
    std::unique_ptr<FILE, decltype(&_pclose)> pipe(
        _popen(("cmd /c " + command).c_str(), "r"), _pclose);
    if (!pipe) {
        throw std::runtime_error("_popen() failed!");
    }

    
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }

    return result;
}



std::string execCommandFromList(std::vector<std::string> Cmd ) {

    std::string command = "";
    for (auto& arg : Cmd) {
        command += arg;
        command += " ";
    }
    std::string result = execCommand(command);
    std::cout << "[EXECCMD] RESULT: " << result << std::endl;

    return result;
}