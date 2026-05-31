package controllers

import (
    "backend/config"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

type TransaksiRequest struct {
    IDProduk          int     `json:"id_produk"`
    Qty               int     `json:"qty"`
    Custcd            string  `json:"custcd"`
    MetodePembayaran  string  `json:"metode_pembayaran"`
    JenisHarga        string  `json:"jenis_harga"`
}

func CreateTransaksi(c *gin.Context) {
    var req TransaksiRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get product price based on jenis_harga
    var hargaProduk float64
    queryHarga := "SELECT harga_produk FROM harga WHERE id_produk = ? AND jenis_harga = ?"
    err := config.DB.QueryRow(queryHarga, req.IDProduk, req.JenisHarga).Scan(&hargaProduk)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product or price type"})
        return
    }

    totalHarga := hargaProduk * float64(req.Qty)

    // Check stock
    var stok int
    queryStok := "SELECT stok_produk FROM produk_ritel WHERE id_produk = ?"
    config.DB.QueryRow(queryStok, req.IDProduk).Scan(&stok)
    
    if stok < req.Qty {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
        return
    }

    // Start transaction
    tx, err := config.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Update stock
    updateStok := "UPDATE produk_ritel SET stok_produk = stok_produk - ? WHERE id_produk = ?"
    _, err = tx.Exec(updateStok, req.Qty, req.IDProduk)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Insert transaction
    query := "INSERT INTO transaksi (id_produk, qty, total_harga, custcd, metode_pembayaran) VALUES (?, ?, ?, ?, ?)"
    result, err := tx.Exec(query, req.IDProduk, req.Qty, totalHarga, req.Custcd, req.MetodePembayaran)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    tx.Commit()
    id, _ := result.LastInsertId()
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Transaction successful",
        "id_transaksi": id,
        "total_harga": totalHarga,
    })
}

func GetTransaksi(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset := (page - 1) * limit
    
    var total int
    config.DB.QueryRow("SELECT COUNT(*) FROM transaksi").Scan(&total)
    
    rows, err := config.DB.Query(`
        SELECT t.id_transaksi, t.id_produk, p.nama_produk, t.qty, t.total_harga, 
               t.custcd, c.nama_customer, t.metode_pembayaran, t.created_at
        FROM transaksi t
        JOIN produk_ritel p ON t.id_produk = p.id_produk
        JOIN customer c ON t.custcd = c.custcd
        ORDER BY t.id_transaksi DESC
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var transaksiList []map[string]interface{}
    for rows.Next() {
        var idTransaksi, idProduk, qty int
        var namaProduk, custcd, namaCustomer, metodePembayaran, createdAt string
        var totalHarga float64
        
        rows.Scan(&idTransaksi, &idProduk, &namaProduk, &qty, &totalHarga, &custcd, &namaCustomer, &metodePembayaran, &createdAt)
        
        transaksi := map[string]interface{}{
            "id_transaksi": idTransaksi,
            "id_produk": idProduk,
            "nama_produk": namaProduk,
            "qty": qty,
            "total_harga": totalHarga,
            "custcd": custcd,
            "nama_customer": namaCustomer,
            "metode_pembayaran": metodePembayaran,
            "created_at": createdAt,
        }
        transaksiList = append(transaksiList, transaksi)
    }

    c.JSON(http.StatusOK, gin.H{
        "data": transaksiList,
        "total": total,
        "page": page,
        "limit": limit,
    })
}

func GetTransaksiByID(c *gin.Context) {
    id := c.Param("id")
    var idTransaksi, idProduk, qty int
    var namaProduk, custcd, namaCustomer, metodePembayaran, createdAt string
    var totalHarga float64
    
    query := `
        SELECT t.id_transaksi, t.id_produk, p.nama_produk, t.qty, t.total_harga, 
               t.custcd, c.nama_customer, t.metode_pembayaran, t.created_at
        FROM transaksi t
        JOIN produk_ritel p ON t.id_produk = p.id_produk
        JOIN customer c ON t.custcd = c.custcd
        WHERE t.id_transaksi = ?
    `
    err := config.DB.QueryRow(query, id).Scan(&idTransaksi, &idProduk, &namaProduk, &qty, &totalHarga, &custcd, &namaCustomer, &metodePembayaran, &createdAt)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "id_transaksi": idTransaksi,
        "id_produk": idProduk,
        "nama_produk": namaProduk,
        "qty": qty,
        "total_harga": totalHarga,
        "custcd": custcd,
        "nama_customer": namaCustomer,
        "metode_pembayaran": metodePembayaran,
        "created_at": createdAt,
    })
}
