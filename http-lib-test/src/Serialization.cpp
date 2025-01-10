#include "../include/Serialization.h"


void from_json(const nlohmann::json& j, ResHeartBeat& hb) {
    j.at("type").get_to(hb.Type);
    j.at("args").get_to(hb.Args);
}

void from_json(const nlohmann::json& j, ResGetSid& obj) {
    j.at("sessionId").get_to(obj.sessionId);
    j.at("welcomeMsg").get_to(obj.welcomeMsg);
}

void from_json(const nlohmann::json& j, ResRegister& obj) {
    j.at("resMsg").get_to(obj.resMsg);
}

void to_json(nlohmann::json& j, const ResHeartBeat& hb) {
    j = nlohmann::json{
        {"type", hb.Type},
        {"args", hb.Args}
    };
}

void to_json(nlohmann::json& j, const ResGetSid& obj) {
    j = nlohmann::json{
        {"sessionId", obj.sessionId},
        {"welcomeMsg", obj.welcomeMsg}
    };
}
