package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/telmomarques/x360h1080p-web-config-server/config"
	"github.com/telmomarques/x360h1080p-web-config-server/customerror"
	"github.com/telmomarques/x360h1080p-web-config-server/hack/rtspserver"
	"github.com/telmomarques/x360h1080p-web-config-server/hack/sshserver"
	"github.com/telmomarques/x360h1080p-web-config-server/hack/websocketstreamserver"
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
	 * RTSP Server
	 */
	rtspServerHackRoutes := apiHackRoutes.Group("/" + rtspserver.ID)

	rtspServerHackRoutes.GET("/config", func(c *gin.Context) {
		c.File(config.GetMetaConfigFilePathForHack(rtspserver.ID))
	})

	rtspServerHackRoutes.GET("/info", func(c *gin.Context) {
		c.String(http.StatusOK, rtspserver.Info())
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

	sshServerHackRoutes.DELETE("/config/users/:username", func(c *gin.Context) {
		var httpStatus = http.StatusOK
		username := c.Param("username")

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
