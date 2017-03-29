package ServerManager

import (
	DownloadServer "UserDownloadServer"
	"golang.org/x/net/context"
	"log"
	"UserServer"
	"errors"
	"database/sql"
)

type UserDownloadServer struct{}

func RunUserDownloadServer() {
	DownloadServer.RegisterUserDownloadServerServer(s, &UserDownloadServer{})
}

func (s *UserDownloadServer) DownloadGame(ctx context.Context, in *DownloadServer.DownloadRequest) (*DownloadServer.DownloadResponse, error) {
	var response DownloadServer.DownloadResponse
	response.Available = false
	response.S3DownloadLink = ""
	
	// Check for valid user credentials
	var login *LoginServer
	loginResponse, err := login.UserLogin(context.Background(), &(UserServer.LoginRequest{Username: in.Username, Password: in.Password}))
	if err != nil || loginResponse.Status != UserServer.LoginReply_OK {
		log.Printf("DownloadGame: invalid credentials provided for user %s", in.Username)
		return &response, errors.New("User login creditials invalid")
	}
	
	var userId string = loginResponse.Id
	
	// check for valid purchase by the user
	var purchaseId string
	err = db.QueryRow("SELECT ID FROM PURCHASE WHERE USER_ID = ? AND GAME_ID = ?", userId, in.GameId).Scan(&purchaseId)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User %s has not purchased game with id %s", in.Username, in.GameId)
		return &response, errors.New("User has not purchased game")
	case err != nil:
		return &response, err
	}
	
	// get S3 link
	err = db.QueryRow("SELECT S3LINK FROM GAME WHERE ID = ?", in.GameId).Scan(&(response.S3DownloadLink))
	if err != nil {
		return &response, err
	}
	
	response.Available = true
	return &response, nil
}

func (s *UserDownloadServer) AvailableDownloads(ctx context.Context, in *DownloadServer.UserDownloadListRequest) (*DownloadServer.DownloadsList, error) {
	var downloads []*DownloadServer.Game = make([]*DownloadServer.Game, 0)
	
	rows, err := db.Query("SELECT GAME_ID FROM PURCHASE WHERE USER_ID = ?", in.Username)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	defer rows.Close()
	for rows.Next() {
		var gameId string
		err := rows.Scan(&gameId)
		if err != nil {
			log.Fatal(err)
		}
		downloads = append(downloads, &DownloadServer.Game{Id : gameId})
	}
	
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	return &DownloadServer.DownloadsList{downloads}, nil
}