package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// InitLogger mengonfigurasi logrus agar menulis ke terminal dan logs/app.log dengan rotasi log otomatis
func InitLogger() {
	Log = logrus.New()

	// 1. Set format ke JSON agar mudah dibaca oleh parser/mesin monitoring
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 2. Buat folder logs jika belum terbentuk
	logDir := "logs"
	_ = os.MkdirAll(logDir, os.ModePerm)

	// 3. Konfigurasi Log Rotation dengan Lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    10,   // Ukuran maksimum berkas log (dalam Megabyte) sebelum dirotasi
		MaxBackups: 3,    // Jumlah maksimum arsip berkas log lama yang disimpan
		MaxAge:     7,    // Masa simpan log lama (dalam hari). Log > 7 hari akan dihapus otomatis
		Compress:   true, // Kompresi berkas log cadangan menjadi format .gz
		LocalTime:  true, // Menggunakan waktu lokal server untuk penulisan timestamp cadangan
	}

	// 4. io.MultiWriter mengalirkan data log ke layar terminal (os.Stdout) dan Lumberjack secara paralel
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	Log.SetOutput(mw)

	// Set level default
	Log.SetLevel(logrus.InfoLevel)
}
