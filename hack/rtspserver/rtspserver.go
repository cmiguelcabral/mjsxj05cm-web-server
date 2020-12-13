package rtspserver

import (
	"os"
	"time"

	"../../network"

	"../../config"
	"../../service"
)

type RTSPServerConfig struct {
	Enable       bool   `json:"enable"`
	EncodingType string `json:"encodingType"`
}

type RTSPServerServiceStatus struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Endpoint string `json:"endpoint"`
}

const ID = "rtsp-server"
const FriendlyName = "Rtsp Server"

const encodingLibPath = "/mnt/data/lib/libboardav.so.1.0.0"
const frameGrabberMainstreamService = "rtsp-server-framegrabber-mainstream"
const frameGrabberSubstreamService = "rtsp-server-framegrabber-substream"
const mainServiceName = "rtsp-server-rtsp-server"

func SaveConfig(rtspServerConfig RTSPServerConfig) bool {
	var h265Byte byte = 0x1
	var h264Byte byte = 0x2

	encodingByte := h265Byte
	if rtspServerConfig.EncodingType == "h264" {
		encodingByte = h264Byte
	}

	_, statError := os.Stat(encodingLibPath)
	if statError != nil {
		return false
	}

	libFile, openError := os.OpenFile(encodingLibPath, os.O_RDWR, 0644)
	if openError != nil {
		return false
	}

	_, writeError := libFile.WriteAt([]byte{encodingByte}, 0x4318)
	if writeError != nil {
		return false
	}

	libFile.Close()

	success := config.Save(ID, rtspServerConfig)

	if !success {
		return false
	}

	if rtspServerConfig.Enable {
		config.EnableService(ID)

		service.Restart(service.Perp, "fetch_av")
		time.Sleep(500 * time.Millisecond)

		service.Restart(service.Runit, frameGrabberMainstreamService)
		service.Restart(service.Runit, frameGrabberSubstreamService)

		service.Restart(service.Runit, mainServiceName)
	} else {
		config.DisableService(ID)

		service.Restart(service.Perp, "fetch_av")

		service.Stop(service.Runit, frameGrabberMainstreamService)
		service.Stop(service.Runit, frameGrabberSubstreamService)

		service.Stop(service.Runit, mainServiceName)
	}

	return true
}

func GetServiceStatus() service.ServiceStatus {
	return service.Status(FriendlyName, service.Runit, mainServiceName)
}

func Info() string {
	ip := network.GetIP()

	return `
		<p>Mainstream (1920x1088): <a href='rtsp://` + ip + `:8554/mainstream'>rtsp://` + ip + `:8554/mainstream</a></p>
		<p>Substream (640x360): <a href='rtsp://` + ip + `:8554/substream'>rtsp://` + ip + `:8554/substream</a></p>
	`
}
