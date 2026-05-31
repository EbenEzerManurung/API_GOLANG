package models

type Harga struct {
    IDProduk     int     `json:"id_produk"`
    NamaProduk   string  `json:"nama_produk"`
    HargaProduk  float64 `json:"harga_produk"`
    JenisHarga   string  `json:"jenis_harga"`
}