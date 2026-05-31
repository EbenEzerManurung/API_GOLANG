package controllers

import (
    "backend/config"
    "backend/models"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset := (page - 1) * limit
    
    var total int
    config.DB.QueryRow("SELECT COUNT(*) FROM customer WHERE inactive_customer = 'N'").Scan(&total)
    
    rows, err := config.DB.Query(`
        SELECT custcd, nama_customer, address, phone, inactive_customer 
        FROM customer 
        WHERE inactive_customer = 'N'
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var customers []models.Customer
    for rows.Next() {
        var customer models.Customer
        rows.Scan(&customer.Custcd, &customer.NamaCustomer, &customer.Address, &customer.Phone, &customer.InactiveCustomer)
        customers = append(customers, customer)
    }

    c.JSON(http.StatusOK, gin.H{
        "data": customers,
        "total": total,
        "page": page,
        "limit": limit,
    })
}

func GetCustomerByID(c *gin.Context) {
    custcd := c.Param("custcd")
    var customer models.Customer
    query := "SELECT custcd, nama_customer, address, phone FROM customer WHERE custcd = ? AND inactive_customer = 'N'"
    err := config.DB.QueryRow(query, custcd).Scan(&customer.Custcd, &customer.NamaCustomer, &customer.Address, &customer.Phone)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
        return
    }
    c.JSON(http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {
    var customer models.Customer
    if err := c.BindJSON(&customer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := "INSERT INTO customer (custcd, nama_customer, address, phone, inactive_customer) VALUES (?, ?, ?, ?, 'N')"
    _, err := config.DB.Exec(query, customer.Custcd, customer.NamaCustomer, customer.Address, customer.Phone)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Customer created successfully"})
}

func UpdateCustomer(c *gin.Context) {
    custcd := c.Param("custcd")
    var customer models.Customer
    if err := c.BindJSON(&customer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := "UPDATE customer SET nama_customer=?, address=?, phone=? WHERE custcd=?"
    _, err := config.DB.Exec(query, customer.NamaCustomer, customer.Address, customer.Phone, custcd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func DeleteCustomer(c *gin.Context) {
    custcd := c.Param("custcd")
    query := "UPDATE customer SET inactive_customer='Y' WHERE custcd=?"
    _, err := config.DB.Exec(query, custcd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
