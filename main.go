package main

import (
	"log"
	"main/ai"
	"main/fungsi"
	"time"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"main/controllers"
	"main/models"

	"main/wa"
	//"main/ai"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// membaca file .env
	godotenv.Load()

	// panggil koneksi
	db := koneksi()

	// auto migrate model ke database
	db.AutoMigrate(&models.Suhu{})
	db.AutoMigrate(&models.Suhu{}, &models.Informasi{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Pesan{})

	// --- TAMBAHKAN INI: Migrate Model Penggajian ---
	db.AutoMigrate(&models.Penggajian{}) //

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// menambahkan middelware JWT
	key_jwt := os.Getenv("KEY_JWT")
	authMiddleware, err := jwtV3.New(&jwtV3.GinJWTMiddleware{
		Realm:       "fikom UDB",
		Key:         []byte(key_jwt),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",

		PayloadFunc: func(data any) jwt.MapClaims {
			value, ok := data.(models.User)
			if ok {
				return jwt.MapClaims{
					"id":   value.ID,
					"nama": value.Nama,
				}
			}
			return jwt.MapClaims{}
		},

		Authenticator: controllers.UserLogin,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// jika route tidak ada
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status": false,
			"pesan":  "Route Tidak Ditemukan",
		})
	})

	// route tanpa middleware
	r.POST("/login", authMiddleware.LoginHandler)

	// route grup dengan middleware jwt
	auth := r.Group("/backend", authMiddleware.MiddlewareFunc())

	auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil",
		})
	})
	auth.POST("/programstudi", fungsi.BacaDataProdi)

	// route suhu
	auth.GET("/suhu", controllers.Tampil)
	auth.POST("/suhu", controllers.Tambah)
	auth.PUT("/suhu", controllers.Ubah)
	auth.DELETE("/suhu", controllers.Hapus)

	// Route Informasi
	auth.GET("/informasi", controllers.InformasiTampil)
	auth.POST("/informasi", controllers.InformasiTambah)
	auth.PUT("/informasi", controllers.InformasiUbah)
	auth.DELETE("/informasi", controllers.InformasiHapus)

	auth.GET("/user", controllers.UserTambah)
	auth.POST("/user", controllers.UserTambah)
	auth.PUT("/user", controllers.UserUbah)
	auth.DELETE("/user", controllers.UserHapus)

	auth.GET("/pesan", controllers.PesanTampil)
	auth.POST("/pesan", controllers.PesanTambah)
	auth.PUT("/pesan", controllers.PesanUbah)
	auth.DELETE("/pesan", controllers.PesanHapus)

	// --- TAMBAHKAN INI: Route CRUD Penggajian ---
	// Endpoint untuk menghitung gaji bulanan karyawan
	auth.POST("/penggajian", controllers.CreatePenggajian)
	auth.GET("/penggajian", controllers.GetPenggajian)
	auth.GET("/penggajian/:id", controllers.GetPenggajianByID)
	auth.PUT("/penggajian/:id", controllers.UpdatePenggajian)
	auth.DELETE("/penggajian/:id", controllers.DeletePenggajian)

	//METHOD DRIVE
	auth.POST("/drive", controllers.DriveUpload)

	// membaca nilai port dari .env
	port := os.Getenv("PORT")
	go r.Run(":" + port)
	//ai.MulaiChatAi()
	ai.InitAi()
	wa.KonekWa(db)

}
