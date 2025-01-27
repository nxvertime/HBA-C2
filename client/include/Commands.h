#ifndef COMMANDS_H
#define COMMANDS_H

#include <string>
#include <vector>

int ExecCommand(std::string command);
int ExecCommandFromList(std::vector<std::string>);

#endif