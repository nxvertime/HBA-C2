#ifndef MESSAGETYPES_H
#define MESSAGETYPES_H

#include <string>
#include <unordered_map>
#include <json.hpp>


struct ResGetSid {
    std::string sessionId;
    std::string welcomeMsg;
};

struct ReqRegister {
    std::string sessionId;
};

struct ReqHeartBeat {
    std::string sessionId;
    std::string type;
    std::string status;
    std::string message;
};

struct ResRegister {
    std::string resMsg;
};

struct ResHeartBeat {
    std::string Type;
    std::unordered_map<std::string, std::string> Args;
};


void from_json(const nlohmann::json& j, ResHeartBeat& hb);
void from_json(const nlohmann::json& j, ResGetSid& obj);
void from_json(const nlohmann::json& j, ResRegister& obj);

void to_json(nlohmann::json& j, const ResHeartBeat& hb);
void to_json(nlohmann::json& j, const ResGetSid& obj);

#endif
