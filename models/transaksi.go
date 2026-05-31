package models

type Transaksi struct {
    IDTransaksi       int     `json:"id_transaksi"`
    IDProduk          int     `json:"id_produk"`
    Qty               int     `json:"qty"`
    TotalHarga        float64 `json:"total_harga"`
    Custcd            string  `json:"custcd"`
    MetodePembayaran  string  `json:"metode_pembayaran"`
}