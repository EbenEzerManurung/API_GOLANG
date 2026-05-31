package routes

import (
    "backend/controllers"
    "backend/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    // Public routes
    r.POST("/api/login", controllers.Login)

    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        api.GET("/users", middleware.AdminOnly(), controllers.GetUsers)
        
        // Customer routes
        api.GET("/customers", controllers.GetCustomers)
        api.GET("/customers/:custcd", controllers.GetCustomerByID)
        api.POST("/customers", middleware.AdminOnly(), controllers.CreateCustomer)
        api.PUT("/customers/:custcd", middleware.AdminOnly(), controllers.UpdateCustomer)
        api.DELETE("/customers/:custcd", middleware.AdminOnly(), controllers.DeleteCustomer)
        
        // Product routes
        api.GET("/produk", controllers.GetProduk)
        api.GET("/produk/:id", controllers.GetProdukByID)
        api.POST("/produk", middleware.AdminOnly(), controllers.CreateProduk)
        api.PUT("/produk/:id", middleware.AdminOnly(), controllers.UpdateProduk)
        api.DELETE("/produk/:id", middleware.AdminOnly(), controllers.DeleteProduk)
        
        // Price routes
        api.GET("/harga", controllers.GetHarga)
        
        // Transaction routes
        api.POST("/transaksi", controllers.CreateTransaksi)
        api.GET("/transaksi", controllers.GetTransaksi)
        api.GET("/transaksi/:id", controllers.GetTransaksiByID)
    }
}
