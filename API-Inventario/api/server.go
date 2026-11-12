package api

/*
Puse algunas importaciones comentadas porque no se están usando, pero las dejo ahí
por si las necesitamos después.
Vamos a necesitarlas despues, pero solo las descomentamos cuando las necesitemos,
porque go no compila si hay importaciones que no se usan.
De momento para hacer los handlers y probarlos con postman,
no las necesitamos, pero luego si vamos a necesitar.
Y con ese codigo es suficiente para hacer los handlers y luego hacer la URL.
*/

import (
	"context"
	"database/sql"

	"rest/repository"
	"rest/security"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	cors "github.com/itsjamie/gin-cors"
)

type Server struct {
	conn                     *sql.DB
	userRepository           *repository.UserRepository
	roleRepository           *repository.RoleRepository
	categoryRepository       *repository.CategoryRepository
	supplyRepository         *repository.SupplyRepository
	supplierRepository       *repository.SupplierRepository
	clientRepository         *repository.ClientRepository
	purchaseRepository       *repository.PurchaseRepository
	purchaseDetailRepository *repository.PurchaseDetailRepository
	saleRepository           *repository.SaleRepository
	saleDetailRepository     *repository.SaleDetailRepository
	reportRepository         *repository.ReportRepository
	backupDir                string
	router                   *gin.Engine
	tokenBuilder             security.Builder
}

func NewServer(
	userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	categoryRepository *repository.CategoryRepository,
	supplyRepository *repository.SupplyRepository,
	supplierRepository *repository.SupplierRepository,
	clientRepository *repository.ClientRepository,
	purchaseRepository *repository.PurchaseRepository,
	purchaseDetailRepository *repository.PurchaseDetailRepository,
	saleRepository *repository.SaleRepository,
	saleDetailRepository *repository.SaleDetailRepository,
	conn *sql.DB, secret string, backupDir string) (*Server, error) {

	builder, err := security.NewPasetoBuilder(secret)
	if err != nil {
		return nil, err
	}

	server := &Server{
		conn:                     conn,
		userRepository:           userRepository,
		roleRepository:           roleRepository,
		categoryRepository:       categoryRepository,
		supplyRepository:         supplyRepository,
		supplierRepository:       supplierRepository,
		clientRepository:         clientRepository,
		purchaseRepository:       purchaseRepository,
		purchaseDetailRepository: purchaseDetailRepository,
		saleRepository:           saleRepository,
		saleDetailRepository:     saleDetailRepository,
		reportRepository:         repository.NewReportRepository(conn),
		backupDir:                backupDir,
		tokenBuilder:             builder,
	}
	router := gin.Default()
	// No confiar en cabeceras X-Forwarded-* de proxies no configurados.
	if err := router.SetTrustedProxies([]string{}); err != nil {
		return nil, err
	}

	router.GET("/health", func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()
		if err := server.conn.PingContext(checkCtx); err != nil {
			ctx.JSON(503, gin.H{"status": "unavailable", "database": "disconnected"})
			return
		}
		ctx.JSON(200, gin.H{"status": "ok", "database": "connected"})
	})

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST,PUT,DELETE",
		RequestHeaders:  "Origin,Authorization,Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	//RUTAS SIN MIDDLEWARE
	router.POST("api/v1/login", server.login)

	//RUTAS CON MIDDLEWARE
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenBuilder))

	//ROUTES ROLES
	authRoutes.GET("api/v1/roles", server.getAllRoles)
	authRoutes.GET("api/v1/roles/:role_id", server.getRoleByID)

	authRoutes.POST("api/v1/roles", server.createRole)
	authRoutes.PUT("api/v1/roles", server.updateRole)
	authRoutes.DELETE("api/v1/roles", server.deleteRoleJSON)
	authRoutes.DELETE("api/v1/roles/:role_id", server.deleteRoleURI)

	//ROUTES CATEGORIES
	authRoutes.GET("api/v1/categories", server.getAllCategories)
	authRoutes.GET("api/v1/categories/:category_id", server.getCategoryByID)

	authRoutes.POST("api/v1/categories", server.createCategory)
	authRoutes.PUT("api/v1/categories", server.updateCategory)
	authRoutes.DELETE("api/v1/categories", server.deleteCategoryJSON)
	authRoutes.DELETE("api/v1/categories/:category_id", server.deleteCategoryURI)

	// ROUTES SUPPLIES
	authRoutes.GET("api/v1/supplies", server.getAllSupplies)
	authRoutes.GET("api/v1/supplies/:supply_id", server.getSupplyByID)

	authRoutes.POST("api/v1/supplies", server.createSupply)
	authRoutes.PUT("api/v1/supplies", server.updateSupply)
	authRoutes.DELETE("api/v1/supplies", server.deleteSupplyJSON)
	authRoutes.DELETE("api/v1/supplies/:supply_id", server.deleteSupplyURI)

	authRoutes.POST("api/v1/supplies/upload_image", server.uploadSupplyImg)
	authRoutes.GET("api/v1/supplies/download_image/:filename", server.downloadSupplyImg)
	authRoutes.DELETE("api/v1/supplies/delete_image/:filename", server.deleteSupplyImg)

	//ROUTES SUPPLIERS
	authRoutes.GET("api/v1/suppliers", server.getAllSuppliers)
	authRoutes.GET("api/v1/suppliers/:supplier_id", server.getSuppliersByID)

	authRoutes.POST("api/v1/suppliers", server.createSupplier)
	authRoutes.PUT("api/v1/suppliers", server.updateSupplier)
	authRoutes.DELETE("api/v1/suppliers", server.deleteSupplierJSON)
	authRoutes.DELETE("api/v1/suppliers/:supplier_id", server.deleteSupplierURI)

	//ROUTES CLIENTS
	authRoutes.GET("api/v1/clients", server.getAllClients)
	authRoutes.GET("api/v1/clients/:client_id", server.getClientByID)

	authRoutes.POST("api/v1/clients", server.createClient)
	authRoutes.PUT("api/v1/clients", server.updateClient)
	authRoutes.DELETE("api/v1/clients", server.deleteClientJSON)
	authRoutes.DELETE("api/v1/clients/:client_id", server.deleteClientURI)

	//ROUTES PURCHASES
	authRoutes.GET("api/v1/purchases", server.getAllPurchases)
	authRoutes.GET("api/v1/purchases/:purchase_id", server.getPurchaseByID)

	authRoutes.POST("api/v1/purchases", server.createPurchase)
	authRoutes.PUT("api/v1/purchases", server.updatePurchase)
	authRoutes.DELETE("api/v1/purchases", server.deletePurchaseJSON)
	authRoutes.DELETE("api/v1/purchases/:purchase_id", server.deletePurchaseURI)

	//ROUTES PURCHASE DETAILS
	authRoutes.GET("api/v1/purchases/:purchase_id/details", server.getPurchaseDetails)
	authRoutes.GET("api/v1/purchase_details/:detail_id", server.getPurchaseDetailByID)

	authRoutes.POST("api/v1/purchase_details", server.createPurchaseDetail)
	authRoutes.PUT("api/v1/purchase_details", server.updatePurchaseDetail)
	authRoutes.DELETE("api/v1/purchase_details", server.deletePurchaseDetailJSON)
	authRoutes.DELETE("api/v1/purchase_details/:detail_id", server.deletePurchaseDetailURI)

	// ROUTES SALES
	authRoutes.GET("api/v1/sales", server.getAllSales)
	authRoutes.GET("api/v1/sales/:sale_id", server.getSaleByID)
	authRoutes.POST("api/v1/sales", server.createSale)
	authRoutes.PUT("api/v1/sales", server.updateSale)
	authRoutes.DELETE("api/v1/sales", server.deleteSaleJSON)
	authRoutes.DELETE("api/v1/sales/:sale_id", server.deleteSaleURI)

	// ROUTES SALE DETAILS
	authRoutes.GET("api/v1/sales/:sale_id/details", server.getSaleDetails)
	authRoutes.GET("api/v1/sale_details/:detail_id", server.getSaleDetailByID)
	authRoutes.POST("api/v1/sale_details", server.createSaleDetail)
	authRoutes.PUT("api/v1/sale_details", server.updateSaleDetail)
	authRoutes.DELETE("api/v1/sale_details", server.deleteSaleDetailJSON)
	authRoutes.DELETE("api/v1/sale_details/:detail_id", server.deleteSaleDetailURI)

	// REPORTES: las consultas consumen vistas SQL camelCase.
	authRoutes.GET("api/v1/reports/sales-by-client", server.getSalesSummaryByClient)
	authRoutes.GET("api/v1/reports/inventory", server.getCurrentInventoryReport)
	authRoutes.GET("api/v1/reports/discount-movements", server.getDiscountMovementsReport)

	// AUDITORÍA Y RESPALDOS (solo administradores).
	authRoutes.GET("api/v1/audit", server.getAuditRecords)
	authRoutes.POST("api/v1/database/backup", server.createDatabaseBackup)
	authRoutes.GET("api/v1/database/backups", server.getDatabaseBackupHistory)

	// ROUTES USERS (administración de usuarios)
	authRoutes.GET("api/v1/users", server.getAllUsers)
	authRoutes.GET("api/v1/users/:user_id", server.getUserByID)
	authRoutes.POST("api/v1/users", server.createUser)
	authRoutes.PUT("api/v1/users", server.updateUser)
	authRoutes.DELETE("api/v1/users", server.deleteUserJSON)
	authRoutes.DELETE("api/v1/users/:user_id", server.deleteUserURI)

	// PERFIL DEL USUARIO AUTENTICADO
	authRoutes.GET("api/v1/profile", server.getProfile)
	authRoutes.PUT("api/v1/profile", server.updateProfile)

	server.router = router
	return server, nil
}

func (server *Server) Start(url string) error {
	return server.router.Run(url)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
