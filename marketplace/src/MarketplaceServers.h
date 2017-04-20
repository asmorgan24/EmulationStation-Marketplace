#ifndef EMULATIONSTATION_ALL_MARKETPLACESERVERS_H
#define EMULATIONSTATION_ALL_MARKETPLACESERVERS_H

#include "GameServer.h"
#include "LoginServer.h"
#include "DownloadServer.h"

class MarketplaceServers {
public:
    static MarketplaceServers* getInstance();
    GameServer* gameServer() { return  gameClient; }
    LoginServer* userServer() { return userClient; }
    DownloadServer* downloadServer() { return downloadClient; }

private:
    MarketplaceServers();
    ~MarketplaceServers();

    static MarketplaceServers* instance;
    GameServer *gameClient;
    LoginServer *userClient;
    DownloadServer *downloadClient;
};


#endif //EMULATIONSTATION_ALL_MARKETPLACESERVERS_H
