package controllers

import (
	"net/http"
	"strconv"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

func ShowUserAddresses(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	var addresses []models.Address
	rows, err := config.DB.Query("SELECT id, address_line1, address_line2, city, state, country, postal_code, is_default FROM address WHERE user_id = $1", userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var address models.Address
		if err := rows.Scan(&address.ID, &address.AddressLine1, &address.AddressLine2, &address.City, &address.State, &address.Country, &address.PostalCode, &address.IsDefault); err != nil {
			c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
			return
		}
		addresses = append(addresses, address)
	}

	c.HTML(http.StatusOK, "address.html", gin.H{"Addresses": addresses})
}

func CreateAddress(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	addressLine1 := c.PostForm("addressLine1")
	addressLine2 := c.PostForm("addressLine2")
	city := c.PostForm("city")
	state := c.PostForm("state")
	country := c.PostForm("country")
	postalCode := c.PostForm("postalCode")

	_, err := config.DB.Exec("INSERT INTO address (user_id, address_line1, address_line2, city, state, country, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7)", userID, addressLine1, addressLine2, city, state, country, postalCode)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/user/address")
}

func SetDefaultAddress(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	addressID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "address.html", gin.H{"Error": "Invalid address ID"})
		return
	}

	// Unset previous default address
	_, err = config.DB.Exec("UPDATE address SET is_default = FALSE WHERE user_id = $1", userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
		return
	}

	// Set new default address
	_, err = config.DB.Exec("UPDATE address SET is_default = TRUE WHERE id = $1 AND user_id = $2", addressID, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/user/address")
}

func DeleteAddress(c *gin.Context) {
	session := utils.GetSession(c.Request)
	userID := session.Values["user_id"]
	if userID == nil {
		c.Redirect(http.StatusFound, "/user/login")
		return
	}

	addressID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "address.html", gin.H{"Error": "Invalid address ID"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM address WHERE id = $1 AND user_id = $2", addressID, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "address.html", gin.H{"Error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/user/address")
}
