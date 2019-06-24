package login

import (
	"net"

	"config"
	"tun"
)
//todo: add sync.Mutx for Users change
type LoginManager struct {
	//key: clientProtocol:clientIP:clientPort  value: key for AES 
	Users map[string]*User
	Tokens map[string]bool
	Cfg *config.
	TunServer *tun.TunServer
	DhcpServer *Dhcp
}

func NewLoginManager(cfg *config.Config, tunServer *tun.TunServer) *LoginManager {
	lm := &LoginManager{
		Users: map[string]*User{},
		Tokens: map[string]bool{},
		Cfg: cfg,
		TunServer: tunServer,
		DhcpServer: NewDhcp(cfg),
	}

	for _, token := range cfg.Tokens {
		lm.Tokens[token] = true
	}
	return lm
}

func (lm *LoginManager) Login(client string, token string) bool {
	if _, ok := lm.Tokens[token]; ok {
		if user, ok := lm.Users[client]; ok {
			user.Close()
		}
		tunAddr, err := lm.DhcpServer.GetNewAddr()
		if err != nil {
			fmt.Println("[LoginManager][Login] no enough ip")
			return false
		}

		user := NewUser(client, tunAddr, token, nil)
		lm.Users[client] = user
		return true
	}
	return false
}

func (lm *LoginManager) Logout(client string) {
	if user, ok := lm.Users[client]; ok {
		user.Close()
		delete(lm.Users, client)
	}
}

func (lm *LoginManager) StartClient(client string, conn net.Conn) {
	if user, ok := lm.Users[client]; ok {
		user.Conn = conn
		user.Start()
		lm.TunServer.StartClient(client, user.OutputChan, user.InputChan)
	}
}

func (lm *LoginManager) GetUser(client string) *login.User{
	if user, ok := lm.Users[client]; ok {
		return user
	}
	return nil
}