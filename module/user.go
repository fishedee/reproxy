package module

import (
	"errors"
	"os"
	"os/user"
	"strconv"
)

type UserConfig struct {
	User string `json:"user"`
}

// 系统用户信息
func InitUser(config UserConfig) error {
	// 配置里的用户信息
	targetUserName := config.User
	if targetUserName == "" {
		return errors.New("invalid user config " + targetUserName)
	}

	// 在系统中寻找该用户
	targetUser, err := user.Lookup(targetUserName)
	if err != nil {
		return err
	}
	targetUserName = targetUser.Username

	// 当前系统用户
	currentUid := os.Getuid()
	currentUser, err := user.LookupId(strconv.Itoa(currentUid))
	if err != nil {
		return err
	}
	currentUserName := currentUser.Username

	// 目标用户与当前用户校验
	if targetUserName != currentUserName {
		return errors.New("invalid current user,target user is " + targetUserName + ",current user is " + currentUserName)
	}

	return nil
}
