package module

import (
	"errors"
	"strconv"
	"os"
	"os/user"
)

type UserConfig struct{
	User string `json:"user"`
}

func InitUser(config UserConfig)(error){
	targetUserName := config.User
	if targetUserName == ""{
		return errors.New("invalid user config "+targetUserName)
	}

	currentUid := os.Getuid()
	targetUser,err := user.Lookup(targetUserName)
	if err != nil{
		return err
	}

	if strconv.Itoa(currentUid) != targetUser.Uid{
		return errors.New("invalid current user,target user is "+targetUserName)
	}

	return nil
}