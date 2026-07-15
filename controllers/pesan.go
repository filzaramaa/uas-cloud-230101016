package controllers

import (
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PesanInput struct {
	Kode    string `json:"kode" binding:"required"`
	Balasan string `json:"balasan" binding:"required"`
}

type UpdatePesanInput struct {
	Balasan string `json:"balasan" binding:"required"`
}

// 1. PesanTambah (CREATE)
func PesanTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input PesanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pesan := models.Pesan{Kode: input.Kode, Balasan: input.Balasan}
	if err := db.Create(&pesan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": pesan})
}

// 2. PesanTampil (READ ALL)
func PesanTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var daftarPesan []models.Pesan
	if err := db.Find(&daftarPesan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": daftarPesan})
}

// 3. PesanUbah (UPDATE)
func PesanUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	kode := c.Param("kode")
	var pesan models.Pesan
	if err := db.First(&pesan, "kode = ?", kode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}
	var input UpdatePesanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&pesan).Updates(models.Pesan{
		Balasan:   input.Balasan,
		UpdatedAt: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{"data": pesan, "message": "Pesan berhasil diperbarui"})
}

// 4. PesanHapus (DELETE)
func PesanHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	kode := c.Param("kode")
	var pesan models.Pesan
	if err := db.First(&pesan, "kode = ?", kode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}
	if err := db.Delete(&pesan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pesan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pesan berhasil dihapus"})
}
