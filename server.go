package main

import (
	"errors"
	"net/http"
	"os"

	"./config"
	"./customerror"
	"./hack/ledcontrol"
	"./hack/motorcontrol"
	"./hack/rtspserver"
	"./hack/sshserver"
	"./hack/websocketstreamserver"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var wwwPath = "/mnt/sdcard/hacks/web-config/www"

var port = "80"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(static.Serve("/js", static.LocalFile(wwwPath+"/js", false)))
	r.Use(static.Serve("/css", static.LocalFile(wwwPath+"/css", false)))

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File(wwwPath + "/favicon.ico")
	})

	r.GET("/", func(c *gin.Context) {
		c.Header("no-store", "expires 0")
		c.File(wwwPath + "/index.html")
	})

	apiHackRoutes := r.Group("/api/hack")

	/**
	 * Led Control
	 */
	ledcontrolHackRoutes := apiHackRoutes.Group("/" + ledcontrol.ID)

	ledcontrolHackRoutes.GET("/state", func(c *gin.Context) {
		led := c.Query("led")
		c.JSON(http.StatusOK, ledcontrol.GetLedStatus(led))
	})

	ledcontrolHackRoutes.GET("/blink", func(c *gin.Context) {
		led := c.Query("led")
		ledcontrol.BlinkLed(led)
		c.String(http.StatusAccepted, "Success!")
	})

	ledcontrolHackRoutes.POST("/state", func(c *gin.Context) {
		var led ledcontrol.Led
		var httpStatus = http.StatusAccepted
		c.Bind(&led)
		success := ledcontrol.SetLed(led)
		if !success {
			httpStatus = http.StatusInternalServerError
		}
		c.String(httpStatus, "Success!")
	})

	/**
	 * Motor Control
	 */
	motorcontrolHackRoutes := apiHackRoutes.Group("/" + motorcontrol.ID)

	motorcontrolHackRoutes.GET("/config", func(c *gin.Context) {
		c.File(config.GetMetaConfigFilePathForHack(motorcontrol.ID))
	})

	motorcontrolHackRoutes.GET("/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, motorcontrol.GetCurrentPosition())
	})

	motorcontrolHackRoutes.POST("/config", func(c *gin.Context) {
		var motorcontrolConfig motorcontrol.MotorControlConfig
		var httpStatus = http.StatusOK
		c.Bind(&motorcontrolConfig)
		success := motorcontrol.SaveConfig(motorcontrolConfig)
		if !success {
			httpStatus = http.StatusInternalServerError
		}
		c.Status(httpStatus)
	})

	motorcontrolHackRoutes.POST("/move", func(c *gin.Context) {
		var motorControlMove motorcontrol.MotorControlMove
		var httpStatus = http.StatusOK
		c.Bind(&motorControlMove)
		success := motorcontrol.MotorMove(motorControlMove)
		if !success {
			httpStatus = http.StatusInternalServerError
		}
		c.Status(httpStatus)
	})

	motorcontrolHackRoutes.POST("/goto", func(c *gin.Context) {
		var motorControlPosition motorcontrol.MotorControlPosition
		var httpStatus = http.StatusOK
		c.Bind(&motorControlPosition)
		success := motorcontrol.MotorGoto(motorControlPosition)
		if !success {
			httpStatus = http.StatusInternalServerError
		}
		c.Status(httpStatus)
	})

	motorcontrolHackRoutes.POST("/command", func(c *gin.Context) {
		var motorControlCommand motorcontrol.MotorControlCommand
		var httpStatus = http.StatusOK
		c.Bind(&motorControlCommand)
		success := motorcontrol.Command(motorControlCommand)
		if !success {
			httpStatus = http.StatusInternalServerError
		}
		c.Status(httpStatus)
	})

	/**
	 * RTSP Server
	 */
	rtspServerHackRoutes := apiHackRoutes.Group("/" + rtspserver.ID)

	rtspServerHackRoutes.GET("/config", func(c *gin.Context) {
		c.File(config.GetMetaConfigFilePathForHack(rtspserver.ID))
	})

	rtspServerHackRoutes.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, rtspserver.Info())
	})

	rtspServerHackRoutes.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, rtspserver.GetServiceStatus())
	})

	rtspServerHackRoutes.POST("/config", func(c *gin.Context) {
		var rtspserverConfig rtspserver.RTSPServerConfig
		var httpStatus = http.StatusOK

		c.Bind(&rtspserverConfig)

		success := rtspserver.SaveConfig(rtspserverConfig)
		if !success {
			httpStatus = http.StatusInternalServerError
		}

		c.Status(httpStatus)
	})

	/**
	 * Websocket Streamer Server
	 */
	websocketStreamerServerHackRoutes := apiHackRoutes.Group("/" + websocketstreamserver.ID)

	websocketStreamerServerHackRoutes.GET("/config", func(c *gin.Context) {
		c.File(config.GetMetaConfigFilePathForHack(websocketstreamserver.ID))
	})

	websocketStreamerServerHackRoutes.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, websocketstreamserver.Info())
	})

	websocketStreamerServerHackRoutes.GET("/endpoints", func(c *gin.Context) {
		c.Data(http.StatusOK, gin.MIMEJSON, []byte(websocketstreamserver.Endpoints()))
	})

	websocketStreamerServerHackRoutes.POST("/config", func(c *gin.Context) {
		var websocketstreamConfig websocketstreamserver.WebsocketStreamConfig
		var httpStatus = http.StatusOK

		c.Bind(&websocketstreamConfig)

		success := websocketstreamserver.SaveConfig(websocketstreamConfig)
		if !success {
			httpStatus = http.StatusInternalServerError
		}

		c.Status(httpStatus)
	})

	/**
	 * SSH/SFTP Server
	 */
	sshServerHackRoutes := apiHackRoutes.Group("/" + sshserver.ID)

	sshServerHackRoutes.GET("/config", func(c *gin.Context) {
		c.File(config.GetMetaConfigFilePathForHack(sshserver.ID))
	})

	sshServerHackRoutes.GET("/config/general", func(c *gin.Context) {
		c.JSON(http.StatusOK, sshserver.GetGeneralConfiguration())
	})

	sshServerHackRoutes.GET("/config/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, sshserver.GetUserConfiguration())
	})

	sshServerHackRoutes.POST("/config/general", func(c *gin.Context) {
		var sshServerConfig sshserver.SSHGeneralConfig
		var httpStatus = http.StatusOK

		c.Bind(&sshServerConfig)

		success := sshserver.SaveGeneralConfig(sshServerConfig)
		if !success {
			httpStatus = http.StatusInternalServerError
		}

		c.Status(httpStatus)
	})

	sshServerHackRoutes.POST("/config/users", func(c *gin.Context) {
		var sshUser sshserver.SSHUser
		var httpStatus = http.StatusOK

		c.Bind(&sshUser)

		err := sshserver.AddUser(sshUser)

		if err != nil {
			var e *customerror.Error

			if errors.As(err, &e) {
				httpStatus = e.HTTPCode
			}
			c.JSON(httpStatus, err)
		}

		c.Status(httpStatus)
	})

	sshServerHackRoutes.DELETE("/config/users", func(c *gin.Context) {
		var httpStatus = http.StatusOK
		username := c.Query("username")

		success := sshserver.DeleteUser(username)
		if !success {
			httpStatus = http.StatusInternalServerError
		}

		c.Status(httpStatus)
	})

	return r
}

func main() {
	if len(os.Args) == 2 {
		wwwPath = os.Args[1]
	}

	r := setupRouter()
	r.Run(":" + port)
}
