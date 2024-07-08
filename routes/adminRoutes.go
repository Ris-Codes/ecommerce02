package routes

import (
	"github.com/Ris-Codes/eCommerce02/controllers"
	middlewares "github.com/Ris-Codes/eCommerce02/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(c *gin.Engine) {

	// Admin Routes
	admin := c.Group("/admin")
	{
		admin.POST("/login", controllers.AdminLogin)
		admin.GET("/login", controllers.ShowAdminLogin)
		admin.GET("/logout", middlewares.AdminAuthRequired(), controllers.AdminLogout)
		admin.GET("/panel", middlewares.AdminAuthRequired(), controllers.AdminPanel)

		// Products
		admin.GET("/panel/products", middlewares.AdminAuthRequired(), controllers.ListOfProducts)
		admin.GET("/panel/products/add_product", middlewares.AdminAuthRequired(), controllers.AddProduct)
		admin.POST("/panel/products/create_product", middlewares.AdminAuthRequired(), controllers.CreateProduct)
		admin.GET("/panel/products/edit_products/:id", middlewares.AdminAuthRequired(), controllers.EditProductPage)
		admin.POST("/panel/products/edit_products/:id", middlewares.AdminAuthRequired(), controllers.EditProduct)
		admin.POST("/panel/products/delete_products/:id", middlewares.AdminAuthRequired(), controllers.DeleteProduct)

		// Categories
		admin.GET("/panel/products/categories", middlewares.AdminAuthRequired(), controllers.ManageCategories)
		admin.POST("/panel/products/category/create", middlewares.AdminAuthRequired(), controllers.CreateCategory)
		admin.POST("/panel/products/category/edit/:id", middlewares.AdminAuthRequired(), controllers.UpdateCategory)
		admin.POST("/panel/products/category/delete/:id", middlewares.AdminAuthRequired(), controllers.DeleteCategory)

		// Users
		admin.GET("/panel/user_management", middlewares.AdminAuthRequired(), controllers.ListOfUsers)
		admin.POST("/panel/user_management/edit_user/:id", middlewares.AdminAuthRequired(), controllers.EditUser)
		admin.POST("/panel/user_management/delete_user/:id", middlewares.AdminAuthRequired(), controllers.DeleteUser)

		// Orders
		admin.GET("/panel/orders", middlewares.AdminAuthRequired(), controllers.ListOfOrders)
		admin.GET("/panel/orders/edit/:id", middlewares.AdminAuthRequired(), controllers.EditOrderPage)
		admin.POST("/panel/orders/edit/:id", middlewares.AdminAuthRequired(), controllers.UpdateOrder)
		admin.POST("/panel/orders/delete/:id", middlewares.AdminAuthRequired(), controllers.DeleteOrder)

		// Super Admin Routes
		admin.GET("/superadmin", middlewares.SuperAdminRequired(), controllers.SuperAdminDashboard)
		admin.POST("/create", middlewares.SuperAdminRequired(), controllers.CreateAdmin)
		admin.GET("/", middlewares.SuperAdminRequired(), controllers.GetAdmins)
		admin.GET("/:id", middlewares.SuperAdminRequired(), controllers.GetAdmin)
		admin.POST("/update/:id", middlewares.SuperAdminRequired(), controllers.UpdateAdmin)
		admin.POST("/delete/:id", middlewares.SuperAdminRequired(), controllers.DeleteAdmin)
	}
}
