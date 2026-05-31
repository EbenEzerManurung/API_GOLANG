package models

type ProdukRitel struct {
    IDProduk       int    `json:"id_produk"`
    NamaProduk     string `json:"nama_produk"`
    StokProduk     int    `json:"stok_produk"`
    InactiveProduk string `json:"inactive_produk"`
}