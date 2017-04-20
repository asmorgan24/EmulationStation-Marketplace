#include "MarketplaceServers.h"

#include <grpc++/create_channel.h>

using namespace std;

MarketplaceServers* MarketplaceServers::instance = NULL;

MarketplaceServers::MarketplaceServers() {
    shared_ptr<grpc::Channel> channel = grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials());
    gameClient = new GameServer(channel);
    downloadClient = new DownloadServer(channel, gameClient);
    userClient = new LoginServer(channel);
}

MarketplaceServers::~MarketplaceServers() {
    if (instance != NULL) {
        delete instance;
        delete gameClient;
        delete downloadClient;
        delete userClient;
    }
}

MarketplaceServers* MarketplaceServers::getInstance() {
    if (instance == NULL) {
        instance = new MarketplaceServers();
    }
    return instance;
}
