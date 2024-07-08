package routes

import (
	"github.com/Ris-Codes/eCommerce02/controllers"
	middlewares "github.com/Ris-Codes/eCommerce02/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine) {

	c.GET("/", middlewares.UserAuthRequired(), controllers.UserHome)
	user := c.Group("/user")
	{
		user.POST("/register", controllers.RegisterUser)
		user.GET("/register", controllers.ShowUserRegister)
		user.POST("/login", controllers.UserLogin)
		user.GET("/login", controllers.ShowUserLogin)

		// Cart
		user.GET("/cart", middlewares.UserAuthRequired(), controllers.Cart)
		user.POST("/cart/add/:id", middlewares.UserAuthRequired(), controllers.AddToCart)
		user.POST("/cart/update/:id", middlewares.UserAuthRequired(), controllers.UpdateCartQuantity)
		user.POST("/cart/remove/:id", middlewares.UserAuthRequired(), controllers.RemoveFromCart)

		// Order and Checkout
		user.POST("/order/create", middlewares.UserAuthRequired(), controllers.CreateOrder)
		user.GET("/payment-success", middlewares.UserAuthRequired(), controllers.PaymentSuccess)
		user.GET("/payment-failure", middlewares.UserAuthRequired(), controllers.PaymentFailure)

		// User Address
		user.GET("/address", middlewares.UserAuthRequired(), controllers.ShowUserAddresses)
		user.POST("address/create", middlewares.UserAuthRequired(), controllers.CreateAddress)
		user.POST("/address/default/:id", middlewares.UserAuthRequired(), controllers.SetDefaultAddress)
		user.POST("/address/delete/:id", middlewares.UserAuthRequired(), controllers.DeleteAddress)

		// User Dashboard
		user.GET("/account", middlewares.UserAuthRequired(), controllers.UserAccount)
		user.POST("/account/update/:id", middlewares.UserAuthRequired(), controllers.UserProfileUpdate)
		user.GET("/logout", middlewares.UserAuthRequired(), controllers.UserLogout)
	}
}
