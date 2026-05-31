package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Konfigurasi database
const (
	DB_USER     = "root"
	DB_PASSWORD = ""
	DB_HOST     = "localhost"
	DB_PORT     = "3306"
	DB_NAME     = "ritel"
)

// Struct untuk menyimpan ID yang baru dibuat
var (
	productIDs   []int
	customerIDs  []string
	userIDs      []int
)

func main() {
	// Koneksi ke database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	// Drop dan buat database baru
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", DB_NAME))
	if err != nil {
		log.Printf("Warning saat drop database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", DB_NAME))
	if err != nil {
		log.Fatal("Gagal membuat database:", err)
	}
	fmt.Printf("✓ Database %s berhasil dibuat (baru & bersih)\n", DB_NAME)

	_, err = db.Exec(fmt.Sprintf("USE %s", DB_NAME))
	if err != nil {
		log.Fatal("Gagal memilih database:", err)
	}

	// Jalankan migration
	fmt.Println("\n📦 Menjalankan migration...")
	runMigration(db)

	fmt.Println("\n🚀 Memulai seeding data...\n")
	startTime := time.Now()

	// Jalankan semua seeder
	seedUsers(db)
	seedCustomers(db, 20)
	seedProducts(db, 50)
	seedPrices(db, 50)
	seedTransactions(db, 100)

	elapsed := time.Since(startTime)

	fmt.Println("\n✅ ========== SEEDING SELESAI ==========")
	fmt.Printf("⏱️  Waktu: %v\n", elapsed)
	fmt.Println("📊 Statistik:")
	fmt.Printf("   - Users: %d\n", len(userIDs))
	fmt.Printf("   - Customers: %d\n", len(customerIDs))
	fmt.Printf("   - Products: %d\n", len(productIDs))
	fmt.Println("=======================================")

	fmt.Println("\n🔐 ========== LOGIN CREDENTIALS ==========")
	fmt.Println("=======================================")
	fmt.Println("| Username       | Password    | Role       |")
	fmt.Println("|----------------|-------------|------------|")
	fmt.Println("| Admin User     | password123 | admin      |")
	fmt.Println("| Kasir User     | password123 | kasir      |")
	fmt.Println("=======================================")
}

// ==================== MIGRATION ====================
func runMigration(db *sql.DB) {
	// Hapus tabel lama jika ada (urutan terbalik karena foreign key)
	tables := []string{"transaksi", "harga", "produk_ritel", "customer", "user"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	}

	// Tabel User (ditambahkan kolom password)
	_, err := db.Exec(`
		CREATE TABLE user (
			id_user INT NOT NULL AUTO_INCREMENT,
			nama_user VARCHAR(100) NOT NULL,
			password VARCHAR(255) NOT NULL,
			role_user ENUM('admin','kasir') NOT NULL,
			PRIMARY KEY (id_user)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel user:", err)
	}
	fmt.Println("   ✓ Migration 1/5 selesai: user")

	// Tabel Customer
	_, err = db.Exec(`
		CREATE TABLE customer (
			custcd VARCHAR(20) NOT NULL,
			nama_customer VARCHAR(200) NOT NULL,
			address TEXT DEFAULT NULL,
			phone VARCHAR(20) DEFAULT NULL,
			inactive_customer ENUM('Y','N') DEFAULT 'N',
			PRIMARY KEY (custcd)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel customer:", err)
	}
	fmt.Println("   ✓ Migration 2/5 selesai: customer")

	// Tabel Produk Ritel
	_, err = db.Exec(`
		CREATE TABLE produk_ritel (
			id_produk INT NOT NULL AUTO_INCREMENT,
			nama_produk VARCHAR(200) NOT NULL,
			stok_produk INT DEFAULT 0,
			inactive_produk ENUM('Y','N') DEFAULT 'N',
			PRIMARY KEY (id_produk)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel produk_ritel:", err)
	}
	fmt.Println("   ✓ Migration 3/5 selesai: produk_ritel")

	// Tabel Harga
	_, err = db.Exec(`
		CREATE TABLE harga (
			id_produk INT NOT NULL,
			nama_produk VARCHAR(200) DEFAULT NULL,
			harga_produk DECIMAL(15,2) DEFAULT NULL,
			jenis_harga ENUM('R','SW','D') NOT NULL,
			PRIMARY KEY (id_produk, jenis_harga),
			FOREIGN KEY (id_produk) REFERENCES produk_ritel(id_produk)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel harga:", err)
	}
	fmt.Println("   ✓ Migration 4/5 selesai: harga")

	// Tabel Transaksi
	_, err = db.Exec(`
		CREATE TABLE transaksi (
			id_transaksi INT NOT NULL AUTO_INCREMENT,
			id_produk INT DEFAULT NULL,
			qty INT NOT NULL,
			total_harga DECIMAL(15,2) DEFAULT NULL,
			custcd VARCHAR(20) DEFAULT NULL,
			metode_pembayaran ENUM('cash','qris','transfer') NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
			PRIMARY KEY (id_transaksi),
			FOREIGN KEY (id_produk) REFERENCES produk_ritel(id_produk),
			FOREIGN KEY (custcd) REFERENCES customer(custcd)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel transaksi:", err)
	}
	fmt.Println("   ✓ Migration 5/5 selesai: transaksi")
}

// ==================== SEEDER USERS ====================
func seedUsers(db *sql.DB) {
	users := []struct {
		namaUser string
		password string
		roleUser string
	}{
		{"Admin User", "password123", "admin"},
		{"Kasir User", "password123", "kasir"},
	}

	fmt.Println("📝 Seeding Users:")
	for _, u := range users {
		result, err := db.Exec(
			"INSERT INTO user (nama_user, password, role_user) VALUES (?, ?, ?)",
			u.namaUser, u.password, u.roleUser,
		)
		if err != nil {
			log.Printf("   ❌ Gagal insert user %s: %v", u.namaUser, err)
			continue
		}

		id, _ := result.LastInsertId()
		userIDs = append(userIDs, int(id))
		fmt.Printf("   ✅ User created: %s (role: %s, password: %s)\n", u.namaUser, u.roleUser, u.password)
	}
	fmt.Printf("✓ Seeder Users: %d users berhasil\n", len(userIDs))
}

// ==================== SEEDER CUSTOMERS ====================
func seedCustomers(db *sql.DB, total int) {
	firstNames := []string{"Budi", "Ani", "Citra", "Dedi", "Eka", "Fajar", "Gina", "Hadi", "Indah", "Joko"}
	lastNames := []string{"Santoso", "Wijaya", "Pratama", "Kurniawan", "Saputra", "Nugroho", "Hidayat"}
	addresses := []string{"Jl. Merdeka No.1", "Jl. Sudirman No.2", "Jl. Thamrin No.3", "Jl. Gatot Subroto No.4"}

	for i := 0; i < total; i++ {
		custcd := fmt.Sprintf("CUST%04d", i+1)
		nama := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
		address := addresses[rand.Intn(len(addresses))]
		phone := fmt.Sprintf("0812%08d", rand.Intn(100000000))
		
		_, err := db.Exec(
			"INSERT INTO customer (custcd, nama_customer, address, phone) VALUES (?, ?, ?, ?)",
			custcd, nama, address, phone,
		)
		if err != nil {
			log.Printf("Gagal insert customer %s: %v", custcd, err)
			continue
		}
		customerIDs = append(customerIDs, custcd)
	}
	fmt.Printf("✓ Seeder Customers: %d customers berhasil\n", len(customerIDs))
}

// ==================== SEEDER PRODUCTS ====================
func seedProducts(db *sql.DB, total int) {
	productNames := []string{
		"Laptop", "Smartphone", "TV LED", "Kulkas", "AC", "Mesin Cuci", 
		"Mouse", "Keyboard", "Monitor", "Printer", "Speaker", "Headset",
		"Kaos", "Kemeja", "Celana Jeans", "Jaket", "Sepatu", "Tas",
		"Roti", "Susu", "Kopi", "Teh", "Gula", "Minyak Goreng",
	}

	for i := 0; i < total; i++ {
		nama := productNames[rand.Intn(len(productNames))] + " " + fmt.Sprintf("V%d", rand.Intn(10)+1)
		stok := rand.Intn(500) + 10
		
		result, err := db.Exec(
			"INSERT INTO produk_ritel (nama_produk, stok_produk) VALUES (?, ?)",
			nama, stok,
		)
		if err != nil {
			log.Printf("Gagal insert product: %v", err)
			continue
		}
		
		id, _ := result.LastInsertId()
		productIDs = append(productIDs, int(id))
	}
	fmt.Printf("✓ Seeder Products: %d products berhasil\n", len(productIDs))
}

// ==================== SEEDER PRICES ====================
func seedPrices(db *sql.DB, total int) {
	jenisHarga := []string{"R", "SW", "D"}
	
	for _, productID := range productIDs {
		// Ambil nama produk
		var namaProduk string
		db.QueryRow("SELECT nama_produk FROM produk_ritel WHERE id_produk = ?", productID).Scan(&namaProduk)
		
		// Untuk setiap jenis harga
		for _, jns := range jenisHarga {
			harga := rand.Intn(500000) + 50000
			_, err := db.Exec(
				"INSERT INTO harga (id_produk, nama_produk, harga_produk, jenis_harga) VALUES (?, ?, ?, ?)",
				productID, namaProduk, harga, jns,
			)
			if err != nil {
				log.Printf("Gagal insert harga for product %d: %v", productID, err)
			}
		}
	}
	fmt.Printf("✓ Seeder Prices: %d price records berhasil\n", len(productIDs)*3)
}

// ==================== SEEDER TRANSACTIONS ====================
func seedTransactions(db *sql.DB, total int) {
	metodeBayar := []string{"cash", "qris", "transfer"}
	
	for i := 0; i < total; i++ {
		if len(productIDs) == 0 || len(customerIDs) == 0 {
			break
		}
		
		productID := productIDs[rand.Intn(len(productIDs))]
		qty := rand.Intn(10) + 1
		
		// Ambil harga dari tabel harga (jenis R = Regular)
		var hargaProduk float64
		db.QueryRow(
			"SELECT harga_produk FROM harga WHERE id_produk = ? AND jenis_harga = 'R'", 
			productID,
		).Scan(&hargaProduk)
		
		totalHarga := float64(qty) * hargaProduk
		custcd := customerIDs[rand.Intn(len(customerIDs))]
		metode := metodeBayar[rand.Intn(len(metodeBayar))]
		
		_, err := db.Exec(
			"INSERT INTO transaksi (id_produk, qty, total_harga, custcd, metode_pembayaran) VALUES (?, ?, ?, ?, ?)",
			productID, qty, totalHarga, custcd, metode,
		)
		if err != nil {
			log.Printf("Gagal insert transaction: %v", err)
		}
	}
	fmt.Printf("✓ Seeder Transactions: %d transactions berhasil\n", total)
}