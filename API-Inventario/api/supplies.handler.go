package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"rest/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Requests para insumos
type createSupplyRequest struct {
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Stock      int32   `json:"stock" binding:"required"`
	CategoryID int32   `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Image      string  `json:"image" binding:"required"`
}

type updateSupplyRequest struct {
	SupplyID   int32   `json:"supply_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Stock      int32   `json:"stock" binding:"required"`
	CategoryID int32   `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Image      string  `json:"image" binding:"required"`
}

type deleteSupplyRequestJSON struct {
	SupplyID int32 `json:"supply_id" binding:"required"`
}

type deleteSupplyRequestURI struct {
	SupplyID int32 `uri:"supply_id" binding:"required"`
}

type getSupplyByIDRequest struct {
	SupplyID int32 `uri:"supply_id" binding:"required"`
}

// Responses para insumos
type categoryResponse struct {
	CategoryID int32  `json:"category_id"`
	Name       string `json:"name"`
}

type SupplyResponse struct {
	SupplyID int32            `json:"supply_id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Stock    int32            `json:"stock"`
	Image    string           `json:"image"`
	Category categoryResponse `json:"category"`
	Price    float64          `json:"price"`
}

type supplyImageRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

// createSupply es el handler para crear un nuevo movimiento.
func (server *Server) createSupply(ctx *gin.Context) {
	var req createSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Transacción con el usuario del token para disAuditoriaInsumos.
	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result, err := repository.NewSupplyRepository(tx).CrearInsumo(
		ctx, req.Name, req.Type, req.Stock, req.CategoryID, req.Price, req.Image)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	lastID := result.SupplyID
	ctx.JSON(http.StatusOK, gin.H{
		"generated_id": lastID,
	})
}

// updateSupply es el handler para actualizar un movimiento.
func (server *Server) updateSupply(ctx *gin.Context) {
	var req updateSupplyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result, err := repository.NewSupplyRepository(tx).ActualizarInsumo(
		ctx, req.SupplyID, req.Name, req.Type, req.Stock, req.CategoryID, req.Price, req.Image)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// deleteSupply elimina un insumo por JSON.
func (server *Server) deleteSupplyJSON(ctx *gin.Context) {
	var req deleteSupplyRequestJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = repository.NewSupplyRepository(tx).EliminarInsumo(ctx, req.SupplyID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Insumo eliminado correctamente",
	})
}

// deleteSupply elimina un insumo por URI.
func (server *Server) deleteSupplyURI(ctx *gin.Context) {
	var req deleteSupplyRequestURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.beginAuditTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = repository.NewSupplyRepository(tx).EliminarInsumo(ctx, req.SupplyID)
	if err = endAuditTx(ctx, tx, err); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Insumo eliminado correctamente",
	})
}

// getAllSupplies obtiene todos los insumos.
// getAllSupplies obtiene todos los insumos.
func (server *Server) getAllSupplies(ctx *gin.Context) {
	fmt.Println("GET /api/v1/supplies: petición recibida")

	// Comprobar si el middleware guardó el usuario autorizado.
	payload, authorized := ctx.Get("authorized")
	fmt.Println("Token autorizado:", authorized)

	if authorized {
		fmt.Printf("Payload autorizado: %+v\n", payload)
	}

	supplies, err := server.supplyRepository.ObtenerInsumos(ctx)
	if err != nil {
		fmt.Printf("ERROR al obtener insumos de la base de datos: %v\n", err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []SupplyResponse

	for _, supply := range supplies {
		response = append(response, SupplyResponse{
			SupplyID: supply.SupplyID,
			Name:     supply.Name,
			Type:     supply.Type,
			Stock:    supply.Stock,
			Image:    supply.Image,
			Price:    supply.Price,
			Category: categoryResponse{
				CategoryID: supply.CategoryID,
				Name:       supply.CategoryName,
			},
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// getSupplyByID obtiene un insumo por su ID.
func (server *Server) getSupplyByID(ctx *gin.Context) {
	var req getSupplyByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supply, err := server.supplyRepository.ObtenerInsumo(ctx, req.SupplyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, supply)
}

/*
// getSuppliesByCategory obtiene los insumos por su categoría.
func (server *Server) getSuppliesByCategory(ctx *gin.Context) {
	var req getSupplyByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	supplies, err := server.dbtx.GetSupplyById(ctx, req.SupplyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	response := categoryResponse{
		CategoryID: supplies.CategoryID,
		Name:       supplies.CategoryName.String,
	}

	ctx.JSON(http.StatusOK, response)
}
*/

/*
func (server *Server) uploadSupplyImg(ctx *gin.Context) {

	fileHeader, err := ctx.FormFile("file0")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer file.Close()

	wd, _ := os.Getwd()
	fmt.Println("Working directory:", wd)
	upDir := filepath.Join(wd, "utils", "data", "supplies")

	err = os.MkdirAll(upDir, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	filename := uuid.New().String() + "_" + filepath.Base(fileHeader.Filename)

	destinationFile, err := os.Create(filepath.Join(upDir, filename))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"message":  "Imagen cargada exitosamente",
	})
}

func (server *Server) downloadSupplyImg(ctx *gin.Context) {

	var req supplyImageRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wd, _ := os.Getwd()

	fileUrl := filepath.Join(
		wd,
		"utils",
		"data",
		"supplies",
		req.Filename,
	)

	ctx.File(fileUrl)
}

func (server *Server) deleteSupplyImg(ctx *gin.Context) {

	var req supplyImageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wd, _ := os.Getwd()

	filePath := filepath.Join(
		wd,
		"utils",
		"data",
		"supplies",
		req.Filename,
	)

	err := os.Remove(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo eliminar la imagen",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Imagen eliminada correctamente",
	})
}
*/

func (server *Server) uploadSupplyImg(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file0")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if fileHeader.Size <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El archivo está vacío"})
		return
	}
	if fileHeader.Size > supplyImageMaxBytes {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "La imagen supera el límite de 2 MB"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, supplyImageMaxBytes+1))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer la imagen"})
		return
	}
	if len(data) == 0 || len(data) > supplyImageMaxBytes {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "La imagen está vacía o supera el límite de 2 MB"})
		return
	}

	contentType := http.DetectContentType(data[:min(len(data), 512)])
	extensionByType := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}
	extension, allowed := extensionByType[contentType]
	if !allowed {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Formato no permitido. Utiliza JPG, PNG o WEBP"})
		return
	}

	filename := uuid.NewString() + extension
	objectPath := "supplies/" + filename
	response, err := supabaseStorageRequest(
		ctx.Request.Context(), http.MethodPost, objectPath,
		bytes.NewReader(data), contentType,
	)
	if err != nil {
		fmt.Printf("ERROR al subir imagen a Supabase Storage: %v\n", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo guardar la imagen en el almacenamiento"})
		return
	}
	response.Body.Close()

	ctx.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"message":  "Imagen cargada exitosamente",
	})
}

func (server *Server) downloadSupplyImg(ctx *gin.Context) {
	var req supplyImageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !validSupplyImageFilename(req.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de imagen no válido"})
		return
	}

	response, err := supabaseStorageRequest(
		ctx.Request.Context(), http.MethodGet, "supplies/"+req.Filename,
		nil, "",
	)
	if err != nil {
		fmt.Printf("ERROR al descargar imagen desde Supabase Storage: %v\n", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo recuperar la imagen del almacenamiento"})
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(io.LimitReader(response.Body, supplyImageMaxBytes+1))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo leer la imagen almacenada"})
		return
	}
	if len(data) > supplyImageMaxBytes {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "La imagen almacenada supera el límite permitido"})
		return
	}

	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data[:min(len(data), 512)])
	}
	ctx.Header("Cache-Control", "private, max-age=300")
	ctx.Data(http.StatusOK, contentType, data)
}

func (server *Server) deleteSupplyImg(ctx *gin.Context) {
	var req supplyImageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !validSupplyImageFilename(req.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de imagen no válido"})
		return
	}

	response, err := supabaseStorageRequest(
		ctx.Request.Context(), http.MethodDelete, "supplies/"+req.Filename,
		nil, "",
	)
	if err != nil {
		fmt.Printf("ERROR al eliminar imagen de Supabase Storage: %v\n", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo eliminar la imagen del almacenamiento"})
		return
	}
	response.Body.Close()

	ctx.JSON(http.StatusOK, gin.H{"message": "Imagen eliminada correctamente"})
}

func validSupplyImageFilename(filename string) bool {
	return filename != "" && filename != "." && filename != ".." &&
		!strings.ContainsAny(filename, "/\\") && filepath.Base(filename) == filename
}
