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
	getAlga(r)
	getAlgae(r)
	updateAlgaConfig(r)
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

func updateAlgaConfig(r *gin.Engine) *gin.Engine {
	r.PUT("/algae/:name/compose", func(c *gin.Context) {
		name := c.Param("name")
		type Body struct {
			Compose string `json:"compose"`
		}
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil || body.Compose == "" {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		logs, err := alga.UpdateAlgaConfig(name, "compose.yml", body.Compose)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "logs": logs})
			return
		}
		c.JSON(200, gin.H{"logs": logs})
	})
	r.PUT("/algae/:name/env", func(c *gin.Context) {
		name := c.Param("name")
		type Body struct {
			Env string `json:"env"`
		}
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil || body.Env == "" {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		logs, err := alga.UpdateAlgaConfig(name, ".env", body.Env)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error(), "logs": logs})
			return
		}
		c.JSON(200, gin.H{"logs": logs})
	})
	return r
}

func getAlga(r *gin.Engine) *gin.Engine {
	r.GET("/algae/:name/", func(c *gin.Context) {
		name := c.Param("name")

		result, err := alga.GetAlga(name)
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"name":    result.Name,
			"compose": result.Compose,
			"env":     result.Env,
		})
	})

	return r
}

func getAlgae(r *gin.Engine) *gin.Engine {
	r.GET("/algae", func(c *gin.Context) {
		result, err := alga.GetAlgae()
		if err != nil {
			c.JSON(err.Code, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
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
