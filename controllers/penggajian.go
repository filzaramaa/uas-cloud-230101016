package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePenggajian menghitung dan menyimpan data gaji sesuai soal UTS
func CreatePenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.Penggajian

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Aturan: Upah lembur per jam adalah Rp 50.000
	const UpahLemburPerJam = 50000.0
	totalUangLembur := float64(input.JamLembur) * UpahLemburPerJam

	// Hitung Gaji Kotor (Gaji Pokok + Total Uang Lembur)
	input.GajiKotor = input.GajiPokok + totalUangLembur

	// Aturan Pajak: Jika Gaji Kotor > Rp 5.000.000, pajak 5%. Jika tidak, 0%[cite: 1]
	if input.GajiKotor > 5000000 {
		input.Pajak = input.GajiKotor * 0.05
	} else {
		input.Pajak = 0
	}

	// Hitung Gaji Bersih (Gaji Kotor - Pajak)[cite: 1]
	input.GajiBersih = input.GajiKotor - input.Pajak

	// Simpan seluruh data ke database[cite: 1]
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data berhasil disimpan", "data": input})
}

// GetPenggajian mengambil semua data penggajian[cite: 1]
func GetPenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var results []models.Penggajian
	db.Find(&results)
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetPenggajianByID mengambil satu data berdasarkan ID[cite: 1]
func GetPenggajianByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var data models.Penggajian
	id := c.Param("id")

	if err := db.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// UpdatePenggajian untuk memperbarui data dan menghitung ulang gaji[cite: 1]
func UpdatePenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var data models.Penggajian
	id := c.Param("id")

	if err := db.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	var input models.Penggajian
	c.ShouldBindJSON(&input)

	// Hitung ulang semua variabel logic[cite: 1]
	const UpahLemburPerJam = 50000.0
	input.GajiKotor = input.GajiPokok + (float64(input.JamLembur) * UpahLemburPerJam)
	if input.GajiKotor > 5000000 {
		input.Pajak = input.GajiKotor * 0.05
	} else {
		input.Pajak = 0
	}
	input.GajiBersih = input.GajiKotor - input.Pajak

	db.Model(&data).Updates(input)
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": data})
}

// DeletePenggajian menghapus data berdasarkan ID[cite: 1]
func DeletePenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var data models.Penggajian
	id := c.Param("id")

	if err := db.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	db.Delete(&data)
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
