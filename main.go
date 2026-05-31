package main

import (
    "backend/config"
    "backend/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    config.InitDB()
    defer config.DB.Close()

    r := gin.Default()
    
    // CORS middleware - allow all origins for Blazor
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Role")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    routes.SetupRoutes(r)
    r.Run(":8080")
}
