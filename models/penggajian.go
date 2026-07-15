package models

import "gorm.io/gorm"

type Penggajian struct {
	gorm.Model
	NamaPegawai string  `json:"nama_pegawai"`
	GajiPokok   float64 `json:"gaji_pokok"`
	JamLembur   int     `json:"jam_lembur"`
	GajiKotor   float64 `json:"gaji_kotor"`
	Pajak       float64 `json:"pajak"`
	GajiBersih  float64 `json:"gaji_bersih"`
}
