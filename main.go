package main

import (
	"fmt"
	"kasir-api/database"
	"kasir-api/handler"
	"kasir-api/repository"
	"kasir-api/service"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBconn string `mapstructure:"DB_CONN"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock int    `json:"stock"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var produk = []Product{
	{ID: 1, Name: "Indomie Goreng", Price: "Rp 3.500", Stock: 100},
	{ID: 2, Name: "Teh Botol Sosro", Price: "Rp 5.000", Stock: 50},
	{ID: 3, Name: "Aqua 600ml", Price: "Rp 4.000", Stock: 200},
}

var category = []Category{
	{ID: 1, Name: "Makanan", Description: "Product makanan ringan dan berat"},
	{ID: 2, Name: "Minuman", Description: "Berbagai jenis minuman kemasan"},
	{ID: 3, Name: "Sembako", Description: "Kebutuhan pokok sehari-hari"},
	{ID: 4, Name: "Snack", Description: "Camilan dan makanan ringan"},
	{ID: 5, Name: "Perlengkapan Rumah", Description: "Alat dan perlengkapan rumah tangga"},
}

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := Config{
		Port:   viper.GetString("PORT"),
		DBconn: viper.GetString("DB_CONN"),
	}

	db, err := database.Connect(config.DBconn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	fmt.Printf("Server running on http://localhost:%s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Error starting server %v\n", err)
	}
}
