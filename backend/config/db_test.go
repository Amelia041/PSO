package config

import (
	"os"
	"testing"
)

// Test 1: GetMongoURI kosong kalau env tidak di-set
func TestGetMongoURIEmpty(t *testing.T) {
	os.Unsetenv("MONGODB_URI")

	uri := GetMongoURI()
	if uri != "" {
		t.Errorf("Harusnya kosong, dapat: %s", uri)
	}
}

// Test 2: GetMongoURI berhasil baca env variable
func TestGetMongoURISet(t *testing.T) {
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGODB_URI")

	uri := GetMongoURI()
	if uri != "mongodb://localhost:27017" {
		t.Errorf("URI tidak sesuai, dapat: %s", uri)
	}
}

// Test 3: Koneksi nyata ke CosmosDB (hanya jalan di CI dengan secret)
func TestConnectDB(t *testing.T) {
	if os.Getenv("MONGODB_URI") == "" {
		t.Skip("MONGODB_URI tidak di-set, skip test koneksi")
	}

	ConnectDB()

	if DB == nil {
		t.Fatal("DB nil setelah ConnectDB — koneksi gagal")
	}
}
