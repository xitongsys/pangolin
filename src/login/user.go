package login

import (
	"encrypt"
)

var CHANBUFFERSIZE = 1024

type User struct {
	Client string
	TunAddr string
	Token string
	Key string
	Chan chan string
}

func NewUser(client string, tun string, token string) *User {
	key := string(encrypt.GetAESKey([]byte(token)))
	return &User {
		Client: client,
		TunAddr: tun,
		Token: token,
		Key: key,
		Chan: make(chan string, CHANBUFFERSIZE),
	}
}