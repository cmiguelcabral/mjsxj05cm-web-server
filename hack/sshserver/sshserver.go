package sshserver

import (
	"../../config"
	"../../customerror"
	"../../service"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
)

const ID = "ssh-server"
const FriendlyName = "SSH / SFTP Server"

const serviceName = "ssh-server-ssh-server"

type SSHUser struct {
	SystemUsername string `json:"systemUsername"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

type SSHServerConfig struct {
	Enable bool      `json:"enable"`
	Users  []SSHUser `json:"users"`
}

type SSHGeneralConfig struct {
	Enable bool `json:"enable"`
}

type SSHUserConfig struct {
	Users []SSHUser `json:"users"`
}

func GetGeneralConfiguration() SSHGeneralConfig {
	var currentConfig SSHServerConfig
	var generalConfig SSHGeneralConfig

	config.Load(ID, &currentConfig)

	generalConfig.Enable = currentConfig.Enable

	return generalConfig
}

func GetUserConfiguration() SSHUserConfig {
	var currentConfig SSHServerConfig
	var userConfig SSHUserConfig

	config.Load(ID, &currentConfig)

	userConfig.Users = currentConfig.Users

	return userConfig
}

func SaveGeneralConfig(newConfig SSHGeneralConfig) bool {
	var updatedconfig SSHServerConfig

	config.Load(ID, &updatedconfig)

	updatedconfig.Enable = newConfig.Enable

	success := config.Save(ID, updatedconfig)

	if !success {
		return false
	}

	if updatedconfig.Enable {
		config.EnableService(ID)
		service.Restart(service.Runit, serviceName)
	} else {
		config.DisableService(ID)
		service.Stop(service.Runit, serviceName)
	}

	return true
}

func AddUser(user SSHUser) error {
	var sshConfig SSHServerConfig

	if user.Username == "" {
		return &customerror.Error{-2, 400, "Username is required."}
	}

	if userExists(user.Username) {
		return &customerror.Error{-3, 400, "Username already exists."}
	}

	if user.Password != "" {
		crypter := sha512_crypt.New()
		salter := sha512_crypt.GetSalt()
		salt := salter.Generate(salter.SaltLenMax)
		passwordHash, _ := crypter.Generate([]byte(user.Password), []byte(salt))
		user.Password = passwordHash
	}

	config.Load(ID, &sshConfig)

	sshConfig.Users = append(sshConfig.Users, user)

	success := config.Save(ID, sshConfig)

	if !success {
		return &customerror.Error{-2, 500, "Error adding user."}
	}

	service.Restart(service.Runit, serviceName)

	return nil
}

func DeleteUser(username string) bool {
	var sshConfig SSHServerConfig

	config.Load(ID, &sshConfig)

	for i, user := range sshConfig.Users {
		if username == user.Username {
			sshConfig.Users = append(sshConfig.Users[:i], sshConfig.Users[i+1:]...)
			break
		}
	}

	success := config.Save(ID, sshConfig)

	if !success {
		return false
	}

	service.Restart(service.Runit, serviceName)
	return true
}

func userExists(username string) bool {
	var sshConfig SSHServerConfig

	config.Load(ID, &sshConfig)

	for _, user := range sshConfig.Users {
		if username == user.Username {
			return true
		}
	}

	return false
}
