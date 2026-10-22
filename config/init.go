package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// var _UPLOAD_DIR string
var _PORT string
var _BASE_DIRECTORY string

func init() {

	// Load env from .env
	// godotenv.Load(".env")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	_BASE_DIRECTORY = dir
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal(err)
	}

	// _UPLOAD_DIR = os.Getenv("UPLOAD_DIR")
	// if _UPLOAD_DIR == "" {
	// 	panic("set UPLOAD_DIR in .env")
	// }
	_PORT = os.Getenv("API_PORT")
	if _PORT == "" {
		_PORT = "8080"
	}
	var devMode bool
	if _devMode := os.Getenv("DEV"); _devMode == "true" {
		devMode = true
	}
	Cfg.api_port = _PORT
	Cfg.base_directory = _BASE_DIRECTORY
	Cfg.dev_mode = devMode
}
