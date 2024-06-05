package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"pharmacy-store/internal/domain/entities"
	"pharmacy-store/internal/infrastructure/messaging/nats"
	"pharmacy-store/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	db         *sql.DB
	natsClient *nats.NatsClient
}

func NewProductHandler(db *sql.DB, natsClient *nats.NatsClient) *ProductHandler {
	return &ProductHandler{db: db, natsClient: natsClient}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []entities.Product
	rows, err := h.db.Query("SELECT id, name, description, price, stock, category_id FROM products")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryID); err != nil {
			utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		products = append(products, product)
	}

	utils.Response(c, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var product entities.Product
	err := h.db.QueryRow("SELECT id, name, description, price, stock, category_id FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Response(c, http.StatusNotFound, "Product Not Found")
			return
		}
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.Response(c, http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err := h.db.QueryRow("INSERT INTO products (name, description, price, stock, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting product into database", "details": err.Error()})
		return
	}

	// Log the product that is being published
	productJSON, _ := json.Marshal(product)
	log.Printf("Publishing product to NATS: %s", string(productJSON))

	h.natsClient.Publish("product.created", product)
	utils.Response(c, http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	_, err := h.db.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4, category_id = $5 WHERE id = $6",
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.natsClient.Publish("product.updated", product)
	utils.Response(c, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := h.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.natsClient.Publish("product.deleted", id)
	utils.Response(c, http.StatusNoContent, nil)
}
