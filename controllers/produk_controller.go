package controllers

import (
    "backend/config"
    "backend/models"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

func GetProduk(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset := (page - 1) * limit
    
    var total int
    config.DB.QueryRow("SELECT COUNT(*) FROM produk_ritel WHERE inactive_produk = 'N'").Scan(&total)
    
    rows, err := config.DB.Query(`
        SELECT id_produk, nama_produk, stok_produk, inactive_produk 
        FROM produk_ritel 
        WHERE inactive_produk = 'N'
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var produkList []models.ProdukRitel
    for rows.Next() {
        var produk models.ProdukRitel
        rows.Scan(&produk.IDProduk, &produk.NamaProduk, &produk.StokProduk, &produk.InactiveProduk)
        produkList = append(produkList, produk)
    }

    c.JSON(http.StatusOK, gin.H{
        "data": produkList,
        "total": total,
        "page": page,
        "limit": limit,
    })
}

func GetProdukByID(c *gin.Context) {
    id := c.Param("id")
    var produk models.ProdukRitel
    query := "SELECT id_produk, nama_produk, stok_produk, inactive_produk FROM produk_ritel WHERE id_produk = ?"
    err := config.DB.QueryRow(query, id).Scan(&produk.IDProduk, &produk.NamaProduk, &produk.StokProduk, &produk.InactiveProduk)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    c.JSON(http.StatusOK, produk)
}

func CreateProduk(c *gin.Context) {
    var produk models.ProdukRitel
    if err := c.BindJSON(&produk); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := "INSERT INTO produk_ritel (nama_produk, stok_produk, inactive_produk) VALUES (?, ?, 'N')"
    result, err := config.DB.Exec(query, produk.NamaProduk, produk.StokProduk)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    id, _ := result.LastInsertId()
    
    // Insert harga dasar
    _, err = config.DB.Exec(`
        INSERT INTO harga (id_produk, nama_produk, harga_produk, jenis_harga) 
        VALUES (?, ?, ?, 'R')
    `, id, produk.NamaProduk, 50000)
    
    c.JSON(http.StatusOK, gin.H{"message": "Product created", "id": id})
}

func UpdateProduk(c *gin.Context) {
    id := c.Param("id")
    var produk models.ProdukRitel
    if err := c.BindJSON(&produk); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := "UPDATE produk_ritel SET nama_produk=?, stok_produk=? WHERE id_produk=?"
    _, err := config.DB.Exec(query, produk.NamaProduk, produk.StokProduk, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProduk(c *gin.Context) {
    id := c.Param("id")
    query := "UPDATE produk_ritel SET inactive_produk='Y' WHERE id_produk=?"
    _, err := config.DB.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
