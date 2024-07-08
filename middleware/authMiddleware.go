package middlewares

import (
	"net/http"

	"github.com/Ris-Codes/eCommerce02/utils"
	"github.com/gin-gonic/gin"
)

func UserAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := utils.GetSession(c.Request)
		if session.Values["user_id"] == nil {
			c.Redirect(http.StatusFound, "/user/login")
			c.Abort()
			return
		}
		if session.IsNew || session.Options.MaxAge <= 0 {
            c.Redirect(http.StatusFound, "/user/login")
            c.Abort()
            return
        }
        c.Next()
	}
}

func AdminAuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := utils.GetSession(c.Request)
        if session.Values["admin_id"] == nil {
            c.Redirect(http.StatusFound, "/admin/login")
            c.Abort()
            return
        }
        if session.IsNew || session.Options.MaxAge <= 0 {
            c.Redirect(http.StatusFound, "/admin/login")
            c.Abort()
            return
        }
        c.Next()
    }
}

func SuperAdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := utils.GetSession(c.Request)
        role := session.Values["role"]
        if role != "superadmin" {
            c.HTML(http.StatusForbidden, "admin_panel.html", gin.H{"error": "Access denied"})
            c.Abort()
            return
        }
        c.Next()
    }
}