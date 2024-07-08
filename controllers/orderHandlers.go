package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	var defaultAddress models.Address
	err := config.DB.QueryRow("SELECT id, address_line1, address_line2, city, state, country, postal_code FROM address WHERE user_id = $1 AND is_default = TRUE", userID).Scan(&defaultAddress.ID, &defaultAddress.AddressLine1, &defaultAddress.AddressLine2, &defaultAddress.City, &defaultAddress.State, &defaultAddress.Country, &defaultAddress.PostalCode)
	if err != nil {
		log.Println("Error querying default address:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "No default address found"})
		return
	}

	rows, err := config.DB.Query("SELECT c.id, c.quantity, p.id, p.product_name, p.price, p.product_image, p.stock FROM cart c JOIN products p ON c.product_id = p.id WHERE c.user_id = $1", userID)
	if err != nil {
		log.Println("Error querying cart items:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var cartItems []models.CartItems
	var totalAmount float64
	for rows.Next() {
		var cartItem models.CartItems
		var product models.Products
		if err := rows.Scan(&cartItem.ID, &cartItem.Quantity, &product.ID, &product.ProductName, &product.Price, &product.ProductImage, &product.Stock); err != nil {
			log.Println("Error scanning cart items:", err)
			c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
			return
		}
		cartItem.Product = product
		cartItems = append(cartItems, cartItem)
		totalAmount += float64(cartItem.Quantity) * float64(product.Price)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}

	// Generate a temporary order ID
	refNumber := generateOrderRefNumber()
	statusID := 1

	var user models.User
	err = config.DB.QueryRow("SELECT id, username, email, phone FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone)
	if err != nil {
		log.Println("Error querying user:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}

	var orderID int
	err = config.DB.QueryRow("INSERT INTO orders (user_id, order_total, ref_number, status_id, shipping_address_id, user_name, user_email, user_phone, payment_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", userID, totalAmount, refNumber, statusID, defaultAddress.ID, user.Username, user.Email, user.Phone, "Pending").Scan(&orderID)
	if err != nil {
		log.Println("Error inserting order:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}

	for _, item := range cartItems {
		_, err = config.DB.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)", orderID, item.Product.ID, item.Quantity, item.Product.Price)
		if err != nil {
			log.Println("Error inserting order items:", err)
			c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
			return
		}
	}

	// Create payment intent with Stripe
	config.InitStripe()
	paymentIntent, err := config.CreatePaymentIntent(int64(totalAmount*100), "inr", "Order #"+refNumber)
	if err != nil {
		log.Println("Error creating payment intent:", err)
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Error creating payment intent"})
		return
	}

	session.Values["order_success"] = true
    session.Save(c.Request, c.Writer)

	c.HTML(http.StatusOK, "checkout.html", gin.H{
		"RefNumber":            refNumber,
		"TotalAmount":          totalAmount,
		"ClientSecret":         paymentIntent.ClientSecret,
		"OrderID":              orderID,
		"OrderItems":           cartItems,
		"UserName":             user.Username,
		"UserEmail":            user.Email,
		"UserPhone":            user.Phone,
		"Address":              defaultAddress,
		"StripePublishableKey": os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}

func PaymentSuccess(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	orderIDStr := c.Query("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "user_home.html", gin.H{"Error": "Invalid order ID"})
		return
	}

	// Update order status
	_, err = config.DB.Exec("UPDATE orders SET payment_status = $1 WHERE id = $2", "Paid", orderID)
	if err != nil {
		log.Println("Error updating order status:", err)
		c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
		return
	}

	// Retrieve order items
	rows, err := config.DB.Query("SELECT product_id, quantity FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		log.Println("Error querying order items:", err)
		c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var orderItems []models.OrderItems
	for rows.Next() {
		var item models.OrderItems
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			log.Println("Error scanning order items:", err)
			c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
			return
		}
		orderItems = append(orderItems, item)
	}

	// Update product stock
	for _, item := range orderItems {
		_, err = config.DB.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.Product.ID)
		if err != nil {
			log.Println("Error updating product stock:", err)
			c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
			return
		}
	}

	// Clear the cart
	_, err = config.DB.Exec("DELETE FROM cart WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error clearing cart:", err)
		c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func PaymentFailure(c *gin.Context) {
	c.Redirect(http.StatusFound, "/user/cart")
}
