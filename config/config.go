package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const hacksRootPath = "/mnt/sdcard/hacks"

const hacksConfigPath = "config"
const hacksMetaConfigPath = "meta/config"
const hacksMetaConfigFilename = "config.json"

const hacksMetaServicePath = "meta/service"

func GetDirectoryPathForHack(hackID string) string {
	return hacksRootPath + "/" + hackID
}

func GetMetaConfigDirectoryPathForHack(hackID string) string {
	return hacksRootPath + "/" + hackID + "/" + hacksMetaConfigPath
}

func GetConfigDirectoryPathForHack(hackID string) string {
	return hacksRootPath + "/" + hackID + "/" + hacksConfigPath
}

func GetMetaConfigFilePathForHack(hackID string) string {
	return GetMetaConfigDirectoryPathForHack(hackID) + "/" + hacksMetaConfigFilename
}

func GetConfigFilePathForHackAndTemplate(hackID string, templateFileName string) string {
	return GetConfigDirectoryPathForHack(hackID) + "/" + templateFileName[0:strings.LastIndex(templateFileName, ".template")]
}

func GetMetaServiceDirectoryPathForHack(hackID string) string {
	return hacksRootPath + "/" + hackID + "/" + hacksMetaServicePath
}

func EnableService(hackID string) {
	os.OpenFile(GetMetaServiceDirectoryPathForHack(hackID)+"/.enable", os.O_RDONLY|os.O_CREATE, 0644)
	exec.Command("/mnt/sdcard/manu_test/configure_services.sh", hackID).Run()
}

func DisableService(hackID string) {
	os.Remove(GetMetaServiceDirectoryPathForHack(hackID) + "/.enable")
}

func Save(hackID string, configStruct interface{}) bool {
	var success bool

	success = writeMetaConfigFile(hackID, configStruct)
	if !success {
		return false
	}

	success = writeConfigFile(hackID, configStruct)
	if !success {
		return false
	}

	return true
}

func Load(hackID string, configStruct interface{}) {
	file, _ := ioutil.ReadFile(GetMetaConfigFilePathForHack(hackID))
	json.Unmarshal([]byte(file), configStruct)
}

func writeMetaConfigFile(hackID string, configStruct interface{}) bool {
	var json, _ = json.MarshalIndent(configStruct, "", "  ")
	writeFileError := ioutil.WriteFile(GetMetaConfigFilePathForHack(hackID), json, 0644)

	if writeFileError != nil {
		return false
	}

	return true
}

func writeConfigFile(hackID string, configStruct interface{}) bool {
	templateFiles := getTemplateFilesForHack(hackID)

	for _, templateFile := range templateFiles {
		templateContent, _ := ioutil.ReadFile(GetMetaConfigDirectoryPathForHack(hackID) + "/" + templateFile)
		parsedTemplate, parseTemplateError := template.New("configTemplate").Parse(string(templateContent))

		if parseTemplateError != nil {
			fmt.Println(parseTemplateError)
			return false
		}

		configFile, openFileError := os.OpenFile(GetConfigFilePathForHackAndTemplate(hackID, templateFile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

		if openFileError != nil {
			fmt.Println(openFileError)
			return false
		}

		executeError := parsedTemplate.Execute(configFile, configStruct)

		if executeError != nil {
			fmt.Println(executeError)
			return false
		}

	}

	return true
}

func getTemplateFilesForHack(hackID string) []string {
	var templateFiles []string

	fileInfo, readDirError := ioutil.ReadDir(GetMetaConfigDirectoryPathForHack(hackID))

	if readDirError == nil {
		for _, file := range fileInfo {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".template") {
				templateFiles = append(templateFiles, file.Name())
			}
		}
	}

	return templateFiles
}
