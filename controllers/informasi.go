package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Menampilkan semua data informasi
func InformasiTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var info []models.Informasi

	db.Find(&info)
	c.JSON(http.StatusOK, gin.H{"status": true, "data": info})
}

// Menambah data informasi
func InformasiTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.Informasi

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": err.Error()})
		return
	}

	db.Create(&input)
	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tambah data", "data": input})
}

// fungsi untuk mengubah data
func InformasiUbah(c *gin.Context) {
	// 1. Ambil koneksi DB dari middleware (Pastikan baris ini ada!)
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": "Koneksi database hilang"})
		return
	}

	var input struct {
		ID         uint   `json:"id"`
		Judul      string `json:"judul"`
		Konten     string `json:"konten"`
		UrlDokumen string `json:"urldokumen"` // Sesuai permintaan tanpa underscore
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Format data salah"})
		return
	}

	var info models.Informasi
	// 2. Cari data lama
	if err := db.First(&info, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "pesan": "Data tidak ditemukan"})
		return
	}

	// 3. Update data
	// Kita gunakan MAP agar GORM tidak bingung dengan tag JSON di model
	err := db.Model(&info).Updates(map[string]interface{}{
		"judul":      input.Judul,
		"konten":     input.Konten,
		"urldokumen": input.UrlDokumen, // Ini merujuk ke NAMA KOLOM di database (biasanya pakai underscore di DB)
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": "Gagal update"})
		return
	}

	// Update data
	db.Model(&info).Updates(models.Informasi{
		Judul:      input.Judul,
		Konten:     input.Konten,
		UrlDokumen: input.UrlDokumen,
	})

	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil ubah data"})
}

// Menghapus data informasi secara PERMANEN (Hard Delete)
func InformasiHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		ID uint `json:"id"`
	}

	// 1. Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "ID harus disertakan dalam JSON"})
		return
	}

	var info models.Informasi

	// 2. Cari data (Gunakan Unscoped agar data yang sudah terhapus "soft delete" juga bisa ditemukan untuk dihapus permanen)
	if err := db.Unscoped().First(&info, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "pesan": "Data tidak ditemukan di database"})
		return
	}

	// 3. Eksekusi Hapus Permanen
	// Pastikan Unscoped() dipanggil sebelum Delete()
	err := db.Unscoped().Delete(&info).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": "Gagal menghapus data dari sistem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Data berhasil dihapus permanen dari database"})
}
