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
	conn *sql.DB, secret string) (*Server, error) {

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
		tokenBuilder:             builder,
	}
	router := gin.Default()

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

	//ROUTES USERS
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
