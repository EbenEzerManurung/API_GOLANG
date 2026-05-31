package controllers

import (
    "backend/config"
    "backend/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetHarga(c *gin.Context) {
    rows, err := config.DB.Query(`
        SELECT h.id_produk, p.nama_produk, h.harga_produk, h.jenis_harga 
        FROM harga h 
        JOIN produk_ritel p ON h.id_produk = p.id_produk 
        WHERE p.inactive_produk = 'N'
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var hargaList []models.Harga
    for rows.Next() {
        var harga models.Harga
        rows.Scan(&harga.IDProduk, &harga.NamaProduk, &harga.HargaProduk, &harga.JenisHarga)
        hargaList = append(hargaList, harga)
    }

    c.JSON(http.StatusOK, hargaList)
}