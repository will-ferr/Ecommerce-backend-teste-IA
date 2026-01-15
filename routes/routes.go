package routes

import (
	"smart-choice/controllers"
	"smart-choice/middlewares"
	"smart-choice/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// Enhanced health checks
	r.GET("/health", controllers.HealthCheck)
	r.GET("/ready", controllers.ReadinessCheck)
	r.GET("/alive", controllers.LivenessCheck)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)

		twofa := auth.Group("/2fa")
		twofa.Use(middlewares.AuthMiddleware())
		{
			twofa.POST("/generate", controllers.Generate2FA)
			twofa.POST("/validate", controllers.Validate2FA)
		}
	}

	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/payment", controllers.PaymentWebhook)
	}

	seo := r.Group("/seo")
	{
		seo.GET("/product/:id", controllers.GetProductMetaTags)

		seo.GET("/category/:category", func(c *gin.Context) {
			category := c.Param("category")
			metaTags := services.GetCategoryMetaTags(category)
			c.JSON(200, metaTags)
		})

		seo.GET("/home", func(c *gin.Context) {
			metaTags := services.GetHomeMetaTags()
			c.JSON(200, metaTags)
		})
	}

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		products := api.Group("/products")
		{
			products.GET("/", controllers.GetProducts)
			products.GET("/:id", controllers.GetProduct)
			products.POST("/", middlewares.AdminMiddleware(), controllers.CreateProduct)
			products.PUT("/:id", middlewares.AdminMiddleware(), controllers.UpdateProduct)
			products.DELETE("/:id", middlewares.AdminMiddleware(), controllers.DeleteProduct)
		}

		coupons := api.Group("/coupons")
		{
			coupons.POST("/validate", controllers.ValidateCoupon)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middlewares.AdminMiddleware())
		{
			dashboard.GET("/metrics", controllers.GetDashboardMetrics)
		}
	}
}
