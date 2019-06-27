package login

import (
	"config"
)

type LoginManager struct {
	//key: clientProtocol:clientIP:clientPort  value: key for AES 
	Users map[string]*User
	Tokens map[string]bool
	Cfg *config.Config
}

func NewLoginManager(cfg *config.Config) *LoginManager {
	lm := &LoginManager{
		Users: map[string]*User{},
		Tokens: map[string]bool{},
		Cfg: cfg,
	}
	for _, token := range cfg.Tokens {
		lm.Tokens[token] = true
	"config"
}
	return lm
}

func (lm *LoginManager) Login(client string, token string) bool {
	if _, ok := lm.Tokens[token]; ok {
	}
	return false
}