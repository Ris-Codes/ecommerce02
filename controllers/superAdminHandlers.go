package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/models"
	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SuperAdminDashboard(c *gin.Context) {
	session := utils.GetSession(c.Request)
	adminID := session.Values["admin_id"]
	adminName := session.Values["admin_name"]
	role := session.Values["role"]

	if role != "superadmin" {
		c.HTML(http.StatusForbidden, "admin_panel.html", gin.H{"error": "Access denied"})
		return
	}

	success := c.Query("success")
    error := c.Query("error")

	c.HTML(http.StatusOK, "superadmin_dashboard.html", gin.H{
		"AdminID":    adminID,
		"AdminName": adminName,
		"Role":       role,
		"success": success,
        "error":   error,
	})
}

func CreateAdmin(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	name := session.Values["admin_name"]
	if role != "superadmin" {
		c.HTML(http.StatusForbidden, "superadmin_dashboard.html", gin.H{"error": "Access denied"})
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	adminRole := c.PostForm("role")

	if adminRole != "admin" && adminRole != "superadmin" {
        c.HTML(http.StatusBadRequest, "superadmin_dashboard.html", gin.H{"error": "Invalid role"})
        return
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "superadmin_dashboard.html", gin.H{"error": "Error hashing password"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO admin (name, password, email, role) VALUES ($1, $2, $3, $4)", username, hashedPassword, email, adminRole)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "superadmin_dashboard.html", gin.H{"error": "Database error"})
		return
	}

	c.HTML(http.StatusOK, "superadmin_dashboard.html", gin.H{
		"AdminName": name,
		"success": "Admin created successfully",})

}

func GetAdmins(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "superadmin" {
		c.HTML(http.StatusForbidden, "superadmin_dashboard.html", gin.H{"error": "Access denied"})
		return
	}

	rows, err := config.DB.Query("SELECT id, name, role FROM admin WHERE role != 'superadmin'")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "superadmin_dashboard.html", gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var admins []models.Admin
	for rows.Next() {
		var admin models.Admin
		err := rows.Scan(&admin.ID, &admin.Name, &admin.Role)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "superadmin_dashboard.html", gin.H{"error": "Database error"})
			return
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
        c.HTML(http.StatusInternalServerError, "admin.html", gin.H{"error": "Database error"})
        return
    }

    success := c.Query("success")
    error := c.Query("error")

    c.HTML(http.StatusOK, "admin.html", gin.H{
        "admins": admins,
        "success": success,
        "error": error,
    })
}

func GetAdmin(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "superadmin" {
		c.HTML(http.StatusForbidden, "superadmin_dashboard.html", gin.H{"error": "Access denied"})
		return
	}

	id := c.Param("id")
	var admin models.Admin
	err := config.DB.QueryRow("SELECT id, name, email, role FROM admin WHERE id=$1", id).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Role)
	if err == sql.ErrNoRows {
		c.HTML(http.StatusNotFound, "superadmin_dashboard.html", gin.H{"error": "Admin not found"})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "superadmin_dashboard.html", gin.H{"error": "Database error"})
		return
	}

	c.HTML(http.StatusOK, "admin.html", gin.H{"admin": admin})
}

func UpdateAdmin(c *gin.Context) {
	session := utils.GetSession(c.Request)
	role := session.Values["role"]
	if role != "superadmin" {
		c.HTML(http.StatusForbidden, "admin.html", gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin.html", gin.H{"error": "Invalid admin ID"})
		return
	}

	username := c.PostForm("username")
	adminRole := c.PostForm("role")

	if adminRole != "admin" && role != "superadmin" {
		c.HTML(http.StatusBadRequest, "admin.html", gin.H{"error": "Invalid role"})
		return
	}

	_, err = config.DB.Exec("UPDATE admin SET name=$1, role=$2 WHERE id=$3", username, adminRole, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin.html", gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, "/admin?success=Admin updated successfully")
}

func DeleteAdmin(c *gin.Context) {
	session := utils.GetSession(c.Request)
    role := session.Values["role"]
    if role != "superadmin" {
        c.HTML(http.StatusForbidden, "admin.html", gin.H{"error": "Access denied"})
        return
    }

    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.HTML(http.StatusBadRequest, "admin.html", gin.H{"error": "Invalid admin ID"})
        return
    }

    _, err = config.DB.Exec("DELETE FROM admin WHERE id=$1", id)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "admin.html", gin.H{"error": "Database error"})
        return
    }

    c.Redirect(http.StatusFound, "/admin?success=Admin deleted successfully")
}
