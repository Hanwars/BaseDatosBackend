package api

import (
	"net/http"
	"rest/security"

	"github.com/gin-gonic/gin"
)

// Requests para categorias
type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryRequest struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type deleteCategoryRequestJSON struct {
	CategoryID int32 `json:"category_id" binding:"required"`
}

type deleteCategoryRequestURI struct {
	CategoryID int32 `uri:"category_id" binding:"required"`
}

type getCategoryByIDRequest struct {
	CategoryID int32 `uri:"category_id" binding:"required"`
}

// createCategory es el handler para crear una nueva categoria.
func (server *Server) createCategory(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req createCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.categoryRepository.CrearCategoria(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	lastID := result.CategoryID
	c.JSON(http.StatusOK, gin.H{
		"generated_id": lastID,
	})
}

// updateCategory es el handler para actualizar una categoria.
func (server *Server) updateCategory(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req updateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.categoryRepository.ActualizarCategoria(c, req.CategoryID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoria actualizado correctamente",
	})
}

// deleteCategoryJSON elimina una categoria por JSON.
func (server *Server) deleteCategoryJSON(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteCategoryRequestJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.categoryRepository.EliminarCategoria(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoria eliminado correctamente",
	})
}

// deleteCategoryURI elimina una categoria por URI.
func (server *Server) deleteCategoryURI(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req deleteCategoryRequestURI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.categoryRepository.EliminarCategoria(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoria eliminado correctamente",
	})
}

// getAllCategories es el handler para obtener todas las categorias.
func (server *Server) getAllCategories(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	categories, err := server.categoryRepository.ObtenerCategorias(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categories)
}

// getCategoryByID es el handler para obtener una categoria por ID.
func (server *Server) getCategoryByID(c *gin.Context) {
	authPayload := c.MustGet("authorized").(*security.Payload)
	if authPayload.Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Solo administradores",
		})
		return
	}

	var req getCategoryByIDRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.categoryRepository.ObtenerCategoria(c, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, category)
}
