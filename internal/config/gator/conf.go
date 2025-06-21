package gator

import "errors"

type Conf struct {
	URL  string `json:"db_url"`
	User string `json:"username"`
}

func (config *Conf) SetUser(username string) error {
	config.User = username
	if config.User != username {
		return errors.New("Failed to set username")
	}
	return nil
}
