package handlers

import (
	"database/sql"
	"net/http"
	"pharmacy-store/internal/domain/entities"
	"pharmacy-store/internal/infrastructure/messaging/nats"
	"pharmacy-store/internal/middleware"
	"pharmacy-store/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db         *sql.DB
	natsClient *nats.NatsClient
}

func NewUserHandler(db *sql.DB, natsClient *nats.NatsClient) *UserHandler {
	return &UserHandler{db: db, natsClient: natsClient}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []entities.User
	rows, err := h.db.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		users = append(users, user)
	}

	utils.Response(c, http.StatusOK, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user entities.User
	err := h.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Response(c, http.StatusNotFound, "User Not Found")
			return
		}
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.Response(c, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	err = h.db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database", "details": err.Error()})
		return
	}

	h.natsClient.Publish("user.created", user)
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	_, err := h.db.Exec("UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		user.Name, user.Email, user.Password, id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.natsClient.Publish("user.updated", user)
	utils.Response(c, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := h.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.natsClient.Publish("user.deleted", id)
	utils.Response(c, http.StatusNoContent, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var user entities.User
	var foundUser entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err := h.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", user.Email).Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Response(c, http.StatusUnauthorized, "User not found")
			return
		}
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		utils.Response(c, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := middleware.GenerateJWT(user.Email)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Error generating token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims, err := middleware.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user entities.User
	err = h.db.QueryRow("SELECT id, name, email FROM users WHERE email = $1", claims.Email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
