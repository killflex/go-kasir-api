# Kasir API

A simple REST API for managing products in a point-of-sale (kasir) system built with Go.

## Features

- Get all products
- Get product by ID
- Create new product
- Update existing product
- Delete product
- Health check endpoint

## Requirements

- Go 1.16 or higher

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

## Notes

- Data is stored in-memory and will reset when the server restarts
- Product IDs are auto-generated and auto-incremented
- All endpoints return JSON responses
