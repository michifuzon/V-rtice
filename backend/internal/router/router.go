package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/controllers"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/middleware"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
	"gorm.io/gorm"
)

func Configurar(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "https://v-rtice-1.onrender.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// DAOs
	usuarioDAO := dao.NuevoUsuarioDAO(db)
	eventoDAO  := dao.NuevoEventoDAO(db)
	entradaDAO := dao.NuevoEntradaDAO(db)
	favoritoDAO := dao.NuevoFavoritoDAO(db)

	// Services
	authService     := services.NuevoAuthService(usuarioDAO)
	eventoService   := services.NuevoEventoService(eventoDAO, entradaDAO, usuarioDAO)
	entradaService  := services.NuevoEntradaService(entradaDAO, eventoDAO, usuarioDAO, db)
	favoritoService := services.NuevoFavoritoService(favoritoDAO, eventoDAO)

	// Controllers
	authController     := controllers.NuevoAuthController(authService)
	eventoController   := controllers.NuevoEventoController(eventoService)
	entradaController  := controllers.NuevoEntradaController(entradaService)
	favoritoController := controllers.NuevoFavoritoController(favoritoService)
	adminController    := controllers.NuevoAdminController(eventoService)
	clubController     := controllers.NuevoClubController(usuarioDAO)
	
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Registrar)
			auth.POST("/login", authController.Login)
		}

		eventos := api.Group("/eventos")
		{
			eventos.GET("", eventoController.Listar)
			eventos.GET("/:id", eventoController.ObtenerPorID)
		}

		// Rutas protegidas — requieren JWT valido
		protegido := api.Group("")
		protegido.Use(middleware.AutenticacionJWT())
		{
			protegido.POST("/entradas", entradaController.Comprar)
			protegido.GET("/entradas/me", entradaController.MisEntradas)
			protegido.DELETE("/entradas/:id", entradaController.Cancelar)
			protegido.POST("/entradas/:id/transfer", entradaController.Traspasar)

			protegido.GET("/club/estado", clubController.EstadoClub)
			protegido.POST("/club/suscribir", clubController.Suscribirse)
			protegido.DELETE("/club/suscripcion", clubController.CancelarSuscripcion)

			protegido.POST("/favoritos/:eventoId", favoritoController.Agregar)
			protegido.DELETE("/favoritos/:eventoId", favoritoController.Quitar)
			protegido.GET("/favoritos", favoritoController.Listar)

			// Rutas de administrador — requieren JWT + rol admin
			admin := protegido.Group("/admin")
			admin.Use(middleware.AutorizarAdmin())
			{
				admin.POST("/eventos", adminController.CrearEvento)
				admin.PUT("/eventos/:id", adminController.ActualizarEvento)
				admin.DELETE("/eventos/:id", adminController.CancelarEvento)
				admin.GET("/eventos/:id/reporte", adminController.ReporteOcupacion)
			}
		}
	}

	return r
}
