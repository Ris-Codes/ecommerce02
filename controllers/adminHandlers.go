package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

// --------------Admin Login---------------
func AdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var dbAdmin models.Admin
	err := config.DB.QueryRow("SELECT id, name, password, role FROM admin WHERE name=$1", username).Scan(&dbAdmin.ID, &dbAdmin.Name, &dbAdmin.Password, &dbAdmin.Role)
	if err == sql.ErrNoRows {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_login.html", gin.H{"error": "Database error"})
		return
	}

	var passwordMatch bool
	if dbAdmin.Role == "superadmin" {
		err := config.DB.QueryRow("SELECT crypt($1, password) = password FROM admin WHERE name=$2", password, username).Scan(&passwordMatch)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "admin_login.html", gin.H{"error": "Database error"})
			return
		}
	} else {
		passwordMatch = checkPasswordHash(password, dbAdmin.Password)
	}

	if !passwordMatch {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	session := utils.GetSession(c.Request)
	session.Values["admin_id"] = dbAdmin.ID
	session.Values["admin_name"] = dbAdmin.Name
	session.Values["role"] = dbAdmin.Role
	utils.SaveSession(c.Request, c.Writer, session)

	if dbAdmin.Role == "superadmin" {
		c.Redirect(http.StatusFound, "/admin/superadmin")
	} else {
		c.Redirect(http.StatusFound, "/admin/panel")
	}
}

func ShowAdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", nil)
}

// --------------Admin Logout---------------
func AdminLogout(c *gin.Context) {
	session := utils.GetSession(c.Request)
	session.Options.MaxAge = -1
	utils.SaveSession(c.Request, c.Writer, session)

	c.Redirect(http.StatusFound, "/admin/login")
}

// -------------Admin Panel----------------
func AdminPanel(c *gin.Context) {
	session := utils.GetSession(c.Request)
	adminID := session.Values["admin_id"]
	adminName := session.Values["admin_name"]
	role := session.Values["role"]

	success := c.Query("success")
	error := c.Query("error")

	c.HTML(http.StatusOK, "admin_panel.html", gin.H{
		"AdminID":   adminID,
		"AdminName": adminName,
		"Role":      role,
		"success":   success,
		"error":     error,
	})
}

// ------------User Management-------------
func ListOfUsers(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" {
		c.HTML(http.StatusForbidden, "admin_panel.html", gin.H{"error": "Access denied"})
		return
	}

	rows, err := config.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_management.html", gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "user_management.html", gin.H{"error": "database error"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.HTML(http.StatusInternalServerError, "user_management.html", gin.H{"error": "database error"})
		return
	}

	success := c.Query("success")
	error := c.Query("error")

	c.HTML(http.StatusOK, "user_management.html", gin.H{
		"users":   users,
		"success": success,
		"error":   error,
	})
}

func EditUser(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" {
		c.HTML(http.StatusForbidden, "user_management.html", gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "user_management.html", gin.H{"error": "Invalid user ID"})
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")

	_, err = config.DB.Exec("UPDATE users SET username=$1, email=$2 WHERE id=$3", username, email, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_management.html", gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/user_management?success=User updated successfully")
}

func DeleteUser(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" {
		c.HTML(http.StatusForbidden, "user_management.html", gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "user_management.html", gin.H{"error": "Invalid user ID"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_management.html", gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/user_management?success=User deleted successfully")
}

// -------------Products----------------
func ListOfProducts(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	var products []models.Products

	searchQuery := c.Query("search")
	if searchQuery != "" && strings.ToLower(searchQuery) != "all" && strings.ToLower(searchQuery) != "show all" {
		rows, err := config.DB.Query("SELECT p.id, p.product_name, p.price, p.description, p.stock, p.product_image, p.category_id, c.id, c.category_name FROM products p JOIN categories c ON p.category_id = c.id WHERE p.product_name ILIKE '%' || $1 || '%'", searchQuery)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var product models.Products
			if err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.Description, &product.Stock, &product.ProductImage, &product.CategoryID, &product.Category.ID, &product.Category.Name); err != nil {
				c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
				return
			}
			products = append(products, product)
		}

		if err := rows.Err(); err != nil {
			c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
			return
		}
	} else {
		rows, err := config.DB.Query("SELECT p.id, p.product_name, p.price, p.description, p.stock, p.product_image, p.category_id, c.id, c.category_name FROM products p JOIN categories c ON p.category_id = c.id")
		if err != nil {
			c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var product models.Products
			if err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.Description, &product.Stock, &product.ProductImage, &product.CategoryID, &product.Category.ID, &product.Category.Name); err != nil {
				c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
				return
			}
			products = append(products, product)
		}

		if err := rows.Err(); err != nil {
			c.HTML(http.StatusInternalServerError, "products.html", gin.H{"Error": "Database error"})
			return
		}
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "products.html", gin.H{
		"Products": products,
		"Search":   searchQuery,
		"Success":  Success,
		"Error":    Error,
	})
}

func AddProduct(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	rows, err := config.DB.Query("SELECT id, category_name FROM categories")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_product.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			c.HTML(http.StatusInternalServerError, "add_product.html", gin.H{"Error": "Database error"})
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		c.HTML(http.StatusInternalServerError, "add_product.html", gin.H{"Error": "Database error"})
		return
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "add_product.html", gin.H{
		"Categories": categories,
		"Success":    Success,
		"Error":      Error,
	})
}

func CreateProduct(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	image := c.PostForm("image")
	categoryID, err := strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "add_product.html", gin.H{"Error": "Invalid category ID"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO products (product_name, price, description, stock, product_image, category_id) VALUES ($1, $2, $3, $4, $5, $6)", name, price, description, stock, image, categoryID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_product.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products/add_product?Success=Product added successfully")
}

func EditProductPage(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_product.html", gin.H{"Error": "Invalid product ID"})
		return
	}

	var product models.Products
	err = config.DB.QueryRow("SELECT id, product_name, price, description, stock, product_image, category_id FROM products WHERE id=$1", id).Scan(&product.ID, &product.ProductName, &product.Price, &product.Description, &product.Stock, &product.ProductImage, &product.CategoryID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_product.html", gin.H{"Error": "Database error"})
		return
	}

	rows, err := config.DB.Query("SELECT id, category_name FROM categories")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_product.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			c.HTML(http.StatusInternalServerError, "edit_product.html", gin.H{"Error": "Database error"})
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		c.HTML(http.StatusInternalServerError, "edit_product.html", gin.H{"Error": "Database error"})
		return
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "edit_product.html", gin.H{
		"Product":    product,
		"Categories": categories,
		"Success":    Success,
		"Error":      Error,
	})
}

func EditProduct(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_product.html", gin.H{"Error": "Invalid product ID"})
		return
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")
	stock := c.PostForm("stock")
	image := c.PostForm("image")
	categoryID, err := strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_product.html", gin.H{"Error": "Invalid category ID"})
		return
	}

	_, err = config.DB.Exec("UPDATE products SET product_name=$1, price=$2, description=$3, stock=$4, product_image=$5, category_id=$6 WHERE id=$7", name, price, description, stock, image, categoryID, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_product.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products?Success=Product updated successfully")
}

func DeleteProduct(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/admin/panel/products?error=Invalid product ID")
		return
	}

	_, err = config.DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/admin/panel/products?error=Database error")
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products?Success=Product deleted successfully")
}

// Categories
func ManageCategories(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	rows, err := config.DB.Query("SELECT id, category_name FROM categories")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "categories.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			c.HTML(http.StatusInternalServerError, "categories.html", gin.H{"Error": "Database error"})
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		c.HTML(http.StatusInternalServerError, "categories.html", gin.H{"Error": "Database error"})
		return
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"Categories": categories,
		"Success":    Success,
		"Error":      Error,
	})
}

func CreateCategory(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	name := c.PostForm("name")

	_, err := config.DB.Exec("INSERT INTO categories (category_name) VALUES ($1)", name)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "manage_categories.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products/categories?Success=Category added successfully")
}

func UpdateCategory(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "categories.html", gin.H{"Error": "Invalid category ID"})
		return
	}

	name := c.PostForm("name")

	_, err = config.DB.Exec("UPDATE categories SET category_name=$1 WHERE id=$2", name, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "categories.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products/categories?Success=Category updated successfully")
}

func DeleteCategory(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "admin" && role != "superadmin" {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusBadRequest, "/admin/panel/prodiucts/categories?error=Invalid category ID")
		return
	}

	_, err = config.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		c.Redirect(http.StatusInternalServerError, "/admin/panel/products/categories?error=Database error")
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/products/categories?Success=Category deleted successfully")
}

// -------------Orders------------------
func ListOfOrders(c *gin.Context) {
	rows, err := config.DB.Query("SELECT o.id, o.user_name, o.user_email, o.user_phone, o.order_total, o.ref_number, o.status_id, o.payment_status, os.status FROM orders o JOIN order_status os ON o.status_id = os.id")
	if err != nil {
		log.Println("Error querying orders:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var orders []models.Orders
	for rows.Next() {
		var order models.Orders
		if err := rows.Scan(&order.ID, &order.UserName, &order.UserEmail, &order.UserPhone, &order.OrderTotal, &order.RefNumber, &order.OrderStatusID, &order.PaymentStatus, &order.OrderStatus); err != nil {
			log.Println("Error scanning orders:", err)
			c.HTML(http.StatusInternalServerError, "admin_order.html", gin.H{"Error": "Database error"})
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		c.HTML(http.StatusInternalServerError, "admin_order.html", gin.H{"Error": "Database error"})
		return
	}

	c.HTML(http.StatusOK, "admin_order.html", gin.H{"Orders": orders})
}

func EditOrderPage(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Orders
	err := config.DB.QueryRow("SELECT id, user_id, user_name,  order_total, ref_number, status_id, payment_status FROM orders WHERE id = $1", orderID).Scan(&order.ID, &order.UserID, &order.UserName, &order.OrderTotal, &order.RefNumber, &order.OrderStatusID, &order.PaymentStatus)
	if err != nil {
		log.Println("Error querying order:", err)
		c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
		return
	}

	// Fetch all order statuses
	rows, err := config.DB.Query("SELECT id, status FROM order_status")
	if err != nil {
		log.Println("Error querying order statuses:", err)
		c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	var statuses []models.OrderStatus
	for rows.Next() {
		var status models.OrderStatus
		if err := rows.Scan(&status.ID, &status.Status); err != nil {
			log.Println("Error scanning order statuses:", err)
			c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
			return
		}
		statuses = append(statuses, status)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
		return
	}

	c.HTML(http.StatusOK, "edit_order.html", gin.H{"Order": order, "Statuses": statuses})
}

func UpdateOrder(c *gin.Context) {
	orderID := c.Param("id") 

	statusID, err := strconv.Atoi(c.PostForm("status_id"))
	if err != nil {
		log.Println("Invalid status ID:", err)
		c.HTML(http.StatusBadRequest, "edit_order.html", gin.H{"Error": "Invalid status ID"})
		return
	}
	paymentStatus := c.PostForm("payment_status")

	_, err = config.DB.Exec("UPDATE orders SET status_id = $1, payment_status = $2 WHERE id = $3", statusID, paymentStatus, orderID)
	if err != nil {
		log.Println("Error updating order:", err)
		c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
		return
	}

	// Check if the order is completed and payment is done
	if statusID == 2 {
		var paymentStatus string
		err := config.DB.QueryRow("SELECT payment_status FROM orders WHERE id = $1", orderID).Scan(&paymentStatus)
        if err != nil {
            log.Println("Error querying payment status:", err)
            c.HTML(http.StatusInternalServerError, "edit_order.html", gin.H{"Error": "Database error"})
            return
        }

        if paymentStatus == "succeeded" {
            // Update the stock and clear the cart
            updateStockAndClearCart(orderID)
        }
    }

	c.Redirect(http.StatusFound, "/admin/panel/orders")
}

func DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	_, err := config.DB.Exec("DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		log.Println("Error deleting order:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/panel/orders")
}

func updateStockAndClearCart(orderID string) {
    // Fetch the order items and update the stock
    rows, err := config.DB.Query("SELECT product_id, quantity FROM order_items WHERE order_id = $1", orderID)
    if err != nil {
        log.Println("Error querying order items:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var productID, quantity int
        if err := rows.Scan(&productID, &quantity); err != nil {
            log.Println("Error scanning order items:", err)
            return
        }

        // Update the stock
        _, err := config.DB.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", quantity, productID)
        if err != nil {
            log.Println("Error updating product stock:", err)
            return
        }
    }

    if err := rows.Err(); err != nil {
        log.Println("Rows error:", err)
        return
    }

    // Clear the cart
    _, err = config.DB.Exec("DELETE FROM cart WHERE user_id = (SELECT user_id FROM orders WHERE id = $1)", orderID)
    if err != nil {
        log.Println("Error clearing cart:", err)
        return
    }
}
