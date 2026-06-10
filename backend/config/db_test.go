package config

import (
	"os"
	"testing"
)

func TestConnectDB(t *testing.T) {
	// Skip kalau MONGODB_URI tidak di-set (misal di local tanpa env)
	if os.Getenv("MONGODB_URI") == "" {
		t.Skip("MONGODB_URI tidak di-set, skip test koneksi")
	}

	// Tidak boleh panic/fatal = koneksi berhasil
	ConnectDB()

	if DB == nil {
		t.Fatal("DB nil setelah ConnectDB — koneksi gagal")
	}
}
