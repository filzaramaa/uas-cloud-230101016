package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DriveUpload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	file, _ := c.FormFile("file")
	mimeType := file.Header.Get("Content-Type")
	fileOpen, _ := file.Open()
	defer fileOpen.Close()
	fileData, _ := ioutil.ReadAll(fileOpen)
	// Encode ke base64.
	data := base64.StdEncoding.EncodeToString(fileData)
	postBody, _ := json.Marshal(map[string]string{
		"fileName": fileName,
		"mimeType": mimeType,
		"data":     data,
	})
	requestBody := bytes.NewBuffer(postBody)
	// post data
	res, err := http.Post(
		"https://script.google.com/macros/s/AKfycbwfx_I2Zwy7rz4CFaNNEpLN9eYzO24cOBVpuKggZTfG3SpGtzxIguoYTkmTkSA5LT8R/exec",
		"application/json; charset=UTF-8",
		requestBody,
	)
	// apakah ada error?
	if err != nil {
		c.JSON(500, gin.H{
			"kode_error": "ERR-DRIVE",
			"pesan":      "Gagal Upload",
		})
		return
	}
	//baca response data
	hasilBody, _ := ioutil.ReadAll(res.Body)
	hasilString := string(hasilBody)
	//konversi string json to object
	var hasilJson map[string]interface{}
	json.Unmarshal([]byte(hasilString), &hasilJson)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Upload",
		"data":   hasilJson,
	})
	// close response body
	res.Body.Close()
}
