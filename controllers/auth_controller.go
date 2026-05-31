package controllers

import (
    "backend/config"
    "backend/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var loginData struct {
        NamaUser string `json:"nama_user"`
        Password string `json:"password"`
    }

    if err := c.BindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
        return
    }

    var user models.User
    query := "SELECT id_user, nama_user, role_user FROM user WHERE nama_user = ?"
    err := config.DB.QueryRow(query, loginData.NamaUser).Scan(&user.IDUser, &user.NamaUser, &user.RoleUser)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    // For demo, password is "password123" for all users
    if loginData.Password != "password123" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": user,
    })
}

func GetUsers(c *gin.Context) {
    rows, err := config.DB.Query("SELECT id_user, nama_user, role_user FROM user")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        rows.Scan(&user.IDUser, &user.NamaUser, &user.RoleUser)
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}