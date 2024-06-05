package handlers

import (
	"database/sql"
	"net/http"
	"pharmacy-store/internal/domain/entities"
	"pharmacy-store/internal/domain/repositories"
	"pharmacy-store/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo *repositories.CategoryRepository
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{
		repo: repositories.NewCategoryRepository(db),
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.Response(c, http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := h.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Response(c, http.StatusNotFound, "Category Not Found")
			return
		}
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.Response(c, http.StatusOK, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category entities.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.repo.Create(&category); err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.Response(c, http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category entities.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.Response(c, http.StatusBadRequest, "Invalid JSON")
		return
	}
	category.ID = id

	if err := h.repo.Update(&category); err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.Response(c, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.repo.Delete(id); err != nil {
		utils.Response(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.Response(c, http.StatusNoContent, nil)
}
