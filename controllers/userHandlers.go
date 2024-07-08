package controllers

import (

	// "github.com/Ris-Codes/eCommerce01/auth"

	"database/sql"
	"net/http"
	"strings"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

// --------------Register User---------------

func ShowUserRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "user_register.html", nil)
}

func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	phone := c.PostForm("phone")

	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (username, password, email, phone) VALUES ($1, $2, $3, $4)", username, hashedPassword, email, phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/user/login")
}

func ShowUserLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "user_login.html", nil)
}

func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var dbUser models.User
	err := config.DB.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !checkPasswordHash(password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := utils.GetSession(c.Request)
	session.Values["user_id"] = dbUser.ID
	session.Values["username"] = dbUser.Username
	utils.SaveSession(c.Request, c.Writer, session)

	c.Redirect(http.StatusFound, "/")
}

func UserHome(c *gin.Context) {
	session := utils.GetSession(c.Request)
	username := session.Values["username"]

	var products []models.Products

	searchQuery := c.Query("search")
	if searchQuery != "" && strings.ToLower(searchQuery) != "all" && strings.ToLower(searchQuery) != "show all" {
		rows, err := config.DB.Query("SELECT p.id, p.product_name, p.price, p.product_image, p.stock, p.description, c.id, c.category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.product_name ILIKE '%' || $1 || '%'", searchQuery)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var product models.Products
			var category models.Categories
			if err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.ProductImage, &product.Stock, &product.Description, &category.ID, &category.Name); err != nil {
				c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
				return
			}
			product.Category = category
			products = append(products, product)
		}

		if err := rows.Err(); err != nil {
			c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
			return
		}
	} else {
		rows, err := config.DB.Query("SELECT p.id, p.product_name, p.price, p.product_image, p.description, p.stock, c.id, c.category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id")

		if err != nil {
			c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Failed to Fetch from Database"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var product models.Products
			var category models.Categories
			err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.ProductImage, &product.Description, &product.Stock, &category.ID, &category.Name)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "user_home.html", gin.H{"Error": "Database error"})
				return
			}
			product.Category = category
			products = append(products, product)
		}
	}

	Success := c.Query("Success")
	Error := c.Query("Error")
	orderSuccess := session.Values["order_success"]
	if orderSuccess != nil && orderSuccess.(bool) {
		session.Values["order_success"] = false // Clear the success message
		session.Save(c.Request, c.Writer)

		c.HTML(http.StatusOK, "user_home.html", gin.H{
			"OrderSuccess": true,
			"Username":     username,
			"Products":     products,
			"Search":       searchQuery,
			"Success":      Success,
			"Error":        Error,
		})
	} else {
		c.HTML(http.StatusOK, "user_home.html", gin.H{
			"Username": username,
			"Products": products,
			"Search":   searchQuery,
			"Success":  Success,
			"Error":    Error,
		})
	}
}

func UserAccount(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, username, email, phone FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_account.html", gin.H{"Error": "Database error"})
		return
	}

	Success := c.Query("Success")
	Error := c.Query("Error")

	c.HTML(http.StatusOK, "user_account.html", gin.H{
		"User":    user,
		"Success": Success,
		"Error":   Error,
	})
}

func UserProfileUpdate(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")
	phone := c.PostForm("phone")

	_, err := config.DB.Exec("UPDATE users SET username=$1, email=$2, phone=$3 WHERE id=$4", username, email, phone, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_account.html", gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/user/account?Success=User updated successfully")
}

func UserLogout(c *gin.Context) {
	session := utils.GetSession(c.Request)
	session.Options.MaxAge = -1
	utils.SaveSession(c.Request, c.Writer, session)

	c.Redirect(http.StatusFound, "/user/login")
}
