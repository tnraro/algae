package api

import (
	"tnraro/algae/internal/alga"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	postAlgae(r)
	deleteAlga(r)
	patchAlga(r)
	postLogin(r)
	return r
}

type postAlgaeBody struct {
	Name    string `json:"name" binding:"required"`
	Compose string `json:"compose" binding:"required"`
	Env     string `json:"env" binding:"required"`
}

func postAlgae(r *gin.Engine) *gin.Engine {
	r.POST("/algae", func(c *gin.Context) {
		var body postAlgaeBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		name, compose, env := body.Name, body.Compose, body.Env
		result, err := alga.CreateAlga(name, compose, env)
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"logs": result})
	})
	return r
}

func deleteAlga(r *gin.Engine) *gin.Engine {
	r.DELETE("/algae/:name", func(c *gin.Context) {
		name := c.Param("name")
		result, err := alga.DeleteAlga(name)
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"logs": result})
	})
	return r
}

type patchAlgaBody struct {
	Compose string `json:"compose"`
	Env     string `json:"env"`
}

func patchAlga(r *gin.Engine) *gin.Engine {
	r.PATCH("/algae/:name", func(c *gin.Context) {
		name := c.Param("name")
		var body patchAlgaBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		result, err := alga.UpdateAlga(name, body.Compose, body.Env)
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"logs": result})
	})
	return r
}

func postLogin(r *gin.Engine) *gin.Engine {
	r.POST("/login", func(c *gin.Context) {
		type Body struct {
			Registry string `json:"registry" binding:"required"`
			Username string `json:"username" binding:"required"`
			Secret   string `json:"secret" binding:"required"`
		}
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		result, err := alga.Login(body.Registry, body.Username, body.Secret)
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"logs": result})
	})
	return r
}
