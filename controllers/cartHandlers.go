package controllers

import (
	"net/http"
	"strconv"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

// ------------Cart-------------
func AddToCart(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "user_home.html", gin.H{"Error": "Invalid product ID"})
		return
	}

    // Check the product stock
    var stock int
    err = config.DB.QueryRow("SELECT stock FROM products WHERE id = $1", productID).Scan(&stock)
    if err != nil {
        c.Redirect(http.StatusFound, "/")
        return
    }

    if stock < 1 {
        c.Redirect(http.StatusFound, "/?Error=out_of_stock")
        return
    }

	var quantity int
    err = config.DB.QueryRow("SELECT quantity FROM cart WHERE user_id=$1 AND product_id=$2", userID, productID).Scan(&quantity)
    if err == nil {
        // Product already in cart, update the quantity
        _, err = config.DB.Exec("UPDATE cart SET quantity = quantity + 1 WHERE user_id = $1 AND product_id = $2", userID, productID)
    } else {
        // Product not in cart, add new entry
        _, err = config.DB.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, 1)", userID, productID)
    }
    if err != nil {
        c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
        return
    }

	c.Redirect(http.StatusFound, "/?Success=Product added to cart successfully")
}

func Cart(c *gin.Context) {

	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	rows, err := config.DB.Query("SELECT c.id, c.quantity, p.id, p.product_name, p.product_image, p.description, p.price, p.stock FROM cart c JOIN products p ON c.product_id = p.id WHERE c.user_id = $1", userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var cartItems []models.CartItems
	var cartTotal float64
	for rows.Next() {
		var cartItem models.CartItems
		var product models.Products
		if err := rows.Scan(&cartItem.ID, &cartItem.Quantity, &product.ID, &product.ProductName, &product.ProductImage, &product.Description, &product.Price, &product.Stock); err != nil {
			c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
			return
		}
		cartItem.Product = product
		cartItems = append(cartItems, cartItem)
		cartTotal += float64(cartItem.Quantity) * float64(product.Price)
		
	}

	if err := rows.Err(); err != nil {
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
		return
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"CartItems": cartItems,
		"CartTotal": cartTotal,
		"Success":   Success,
		"Error":     Error,
	})
}

func UpdateCartQuantity(c *gin.Context) {
    session := utils.GetSession(c.Request)
    userID := session.Values["user_id"]
    if userID == nil {
        c.Redirect(http.StatusFound, "/user/login")
        return
    }

    cartItemID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.HTML(http.StatusBadRequest, "cart.html", gin.H{"Error": "Invalid cart item ID"})
        return
    }

    action := c.PostForm("action")
    var change int
    if action == "increase" {
        change = 1
    } else if action == "decrease" {
        change = -1
    } else {
        c.HTML(http.StatusBadRequest, "cart.html", gin.H{"Error": "Invalid action"})
        return
    }

    var currentQuantity int
    var productStock int
    err = config.DB.QueryRow("SELECT c.quantity, p.stock FROM cart c JOIN products p ON c.product_id = p.id WHERE c.id = $1", cartItemID).Scan(&currentQuantity, &productStock)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
        return
    }

    newQuantity := currentQuantity + change
    if newQuantity < 1 || newQuantity > productStock {
        c.HTML(http.StatusBadRequest, "cart.html", gin.H{"Error": "Out of Stock"})
        return
    }

    _, err = config.DB.Exec("UPDATE cart SET quantity = $1 WHERE id = $2 AND user_id = $3", newQuantity, cartItemID, userID)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
        return
    }

    c.Redirect(http.StatusFound, "/user/cart?Success=Quantity updated successfully")
}

func RemoveFromCart(c *gin.Context) {
    session := utils.GetSession(c.Request)
    userID := session.Values["user_id"]
    if userID == nil {
        c.Redirect(http.StatusFound, "/user/login")
        return
    }

    cartItemID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.HTML(http.StatusBadRequest, "cart.html", gin.H{"Error": "Invalid cart item ID"})
        return
    }

    _, err = config.DB.Exec("DELETE FROM cart WHERE id=$1 AND user_id=$2", cartItemID, userID)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"Error": "Database error"})
        return
    }

    c.Redirect(http.StatusFound, "/user/cart?Success=Item removed from cart successfully")
}