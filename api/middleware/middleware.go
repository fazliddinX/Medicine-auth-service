package middleware

import (
	"auth-service/pkg/models"
	"auth-service/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAccessTokenMid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Cookie'dan tokenni olish
		tokenString, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.Error{"Authorization cookie not found"})
			ctx.Abort()
			return
		}

		// JWT tokenni tekshirish va tasdiqlash
		claims, err := token.ExtractClaimsRefresh(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.Error{
				Error: "Invalid token",
			})
			ctx.Abort()
			return
		}

		// Foydalanuvchi ma'lumotlarini context ga qo'shish
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
