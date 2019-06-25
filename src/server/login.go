package server

import (
	"net"
	"fmt"
	"sync"

	"config"
	"tun"
	"logging"
)
//todo: add sync.Mutx for Users change
type LoginManager struct {
	//key: clientProtocol:clientIP:clientPort  value: key for AES 
	Users map[string]*User
	Tokens map[string]bool
	Cfg *config.Config
	TunServer *tun.TunServer
	DhcpServer *Dhcp

	Mutex sync.Mutex
}

func NewLoginManager(cfg *config.Config) (*LoginManager, error) {
	tunServer, err := tun.NewTunServer(cfg.TunName, cfg.Mtu)
	if err != nil {
		return nil, err
	}

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
	return lm, nil
}

func (lm *LoginManager) Login(client string, token string) error {
	defer lm.Mutex.Unlock()
	lm.Mutex.Lock()
	if _, ok := lm.Tokens[token]; ok {
		if user, ok := lm.Users[client]; ok {
			user.Close()
		}
		localTunIp, err := lm.DhcpServer.ApplyIp()
		if err != nil {
			return err
		}

		user := NewUser(client, localTunIp, token, nil, lm.Logout)
		lm.Users[client] = user

		logging.Log.Info("User login: client: %v localTunIp: %v", user.Client, user.LocalTunIp)
		return nil
	}
	return fmt.Errorf("token not found")
}

func (lm *LoginManager) Logout(client string) {
	defer lm.Mutex.Unlock()
	lm.Mutex.Lock()
	if user, ok := lm.Users[client]; ok {
		lm.DhcpServer.ReleaseIp(user.LocalTunIp)
		delete(lm.Users, client)

		logging.Log.Info("User logout: client: %v localTunIp: %v", user.Client, user.LocalTunIp)
	}
}

func (lm *LoginManager) Start() {
	lm.TunServer.Start()
}

func (lm *LoginManager) StartClient(client string, conn net.Conn) {
	if user, ok := lm.Users[client]; ok {
		user.Conn = conn
		user.Start()
		lm.TunServer.StartClient(client, user.ConnToTunChan, user.TunToConnChan)
	}
}

func (lm *LoginManager) GetUser(client string) *User{
	if user, ok := lm.Users[client]; ok {
		return user
	}
	return nil
}