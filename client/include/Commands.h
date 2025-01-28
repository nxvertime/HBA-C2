#ifndef COMMANDS_H
#define COMMANDS_H

#include <string>
#include <vector>
#include <Windows.h>

#include <iostream>
#include <cstdio>
#include <memory>

#include <array>

std::string execCommandFromList(std::vector<std::string>);
std::string execCommand(const std::string& command);
#endif