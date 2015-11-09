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
	targetUser,err := user.Lookup(targetUserName)
	if err != nil{
		return err
	}
	targetUserName = targetUser.Username

	currentUid := os.Getuid()
	currentUser,err := user.LookupId(strconv.Itoa(currentUid))
	if err != nil{
		return err
	}
	currentUserName := currentUser.Username

	if targetUserName != currentUserName{
		return errors.New("invalid current user,target user is "+targetUserName+",current user is "+currentUserName)
	}

	return nil
}