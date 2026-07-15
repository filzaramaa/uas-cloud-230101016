package controllers

import (
	"crypto/sha1"
	"fmt"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	jwt "github.com/appleboy/gin-jwt/v3"
)

// ================== STRUCT ==================

type StrukturUserTambah struct {
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StrukturUserUbah struct {
	Id       uint   `json:"id" binding:"required"`
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StrukturUserHapus struct {
	Id uint `json:"id" binding:"required"`
}

type StrukturLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ================== USER TAMBAH ==================

func UserTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserTambah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	// enkripsi password
	sha := sha1.New()
	sha.Write([]byte(dataUser.Password))
	encrypted := sha.Sum(nil)
	encryptedString := fmt.Sprintf("%x", encrypted)

	modelUser := models.User{
		Nama:     dataUser.Nama,
		Username: dataUser.Username,
		Password: encryptedString,
	}

	hasil := db.Create(&modelUser)
	if hasil.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah data",
			"kesalahan": hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil tambah data",
		"data":   modelUser,
	})
}

// ================== USER UBAH ==================

func UserUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserUbah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelUser models.User
	cekUser := db.First(&modelUser, dataUser.Id)

	if cekUser.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "User tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
		})
		return
	}

	// enkripsi password
	sha := sha1.New()
	sha.Write([]byte(dataUser.Password))
	encrypted := sha.Sum(nil)
	encryptedString := fmt.Sprintf("%x", encrypted)

	modelUser.Nama = dataUser.Nama
	modelUser.Username = dataUser.Username
	modelUser.Password = encryptedString

	hasil := db.Save(&modelUser)
	if hasil.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah data",
			"kesalahan": hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil ubah data",
		"data":   modelUser,
	})
}

// ================== USER HAPUS ==================

func UserHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserHapus
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelUser models.User
	hasil := db.Delete(&modelUser, dataUser.Id)

	if hasil.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus data",
			"kesalahan": hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil hapus data",
	})
}

// ================== USER LOGIN (JWT) ==================

func UserLogin(c *gin.Context) (any, error) {
	dbInterface, exists := c.Get("db")
	if !exists {
		return nil, jwt.ErrFailedAuthentication
	}

	db := dbInterface.(*gorm.DB)

	var dataUser StrukturLogin
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	// validasi input kosong
	if dataUser.Username == "" || dataUser.Password == "" {
		return nil, jwt.ErrMissingLoginValues
	}

	// enkripsi password (SHA1)
	sha := sha1.New()
	sha.Write([]byte(dataUser.Password))
	encrypted := sha.Sum(nil)
	encryptedString := fmt.Sprintf("%x", encrypted)

	var modelUser models.User

	cekUser := db.Where("username = ?", dataUser.Username).
		Where("password = ?", encryptedString).
		First(&modelUser)

	// 🔥 handle error lebih jelas
	if cekUser.Error != nil {
		if cekUser.Error == gorm.ErrRecordNotFound {
			return nil, jwt.ErrFailedAuthentication
		}
		// kalau error lain (DB error)
		fmt.Println("DB ERROR:", cekUser.Error)
		return nil, cekUser.Error
	}

	return modelUser, nil
}
