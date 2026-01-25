# Kasir API

A simple REST API for managing products and categories in a point-of-sale (kasir) system built with Go.

## Features

- Health check endpoint
- **Products Management:**
  - Get all products
  - Create new product
  - Get product by ID
  - Update existing product by ID
  - Delete product by ID
- **Category Management:**
  - Get all categories
  - Create new category
  - Get category by ID
  - Update existing category by ID
  - Delete category by ID

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

```
GET /health
```

Response:

```json
{
  "message": "API is running"
}
```

### Products

**Get all products**

```
GET /api/produk
```

**Get product by ID**

```
GET /api/produk/{id}
```

**Create new product**

```
POST /api/produk
Content-Type: application/json

{
  "name": "Product Name",
  "harga": "Rp 10.000",
  "stok": 50
}
```

**Update product**

```
PUT /api/produk/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "harga": "Rp 12.000",
  "stok": 75
}
```

**Delete product**

```
DELETE /api/produk/{id}
```

### Categories

**Get all categories**

```
GET /category
```

**Get category by ID**

```
GET /category/{id}
```

**Create new category**

```
POST /category
Content-Type: application/json

{
  "name": "Category Name",
  "description": "Category description"
}
```

**Update category**

```
PUT /category/{id}
Content-Type: application/json

{
  "name": "Updated Category Name",
  "description": "Updated description"
}
```

**Delete category**

```
DELETE /category/{id}
```

## Response Format

Success responses include a message and data object:

```json
{
  "message": "Operation successful",
  "data": { ... }
}
```

Error responses:

```json
{
  "error": "Error description"
}
```

## Sample Data

### Products

- Indomie Goreng - Rp 3.500 (Stock: 100)
- Teh Botol Sosro - Rp 5.000 (Stock: 50)
- Aqua 600ml - Rp 4.000 (Stock: 200)

### Categories

- Makanan - Produk makanan ringan dan berat
- Minuman - Berbagai jenis minuman kemasan
- Sembako - Kebutuhan pokok sehari-hari
- Snack - Camilan dan makanan ringan
- Perlengkapan Rumah - Alat dan perlengkapan rumah tangga

## Notes

- Data is stored in-memory and will reset when the server restarts
- Product and Category IDs are auto-generated and auto-incremented
- All endpoints return JSON responses
- The server can use a custom port by setting the `PORT` environment variable
