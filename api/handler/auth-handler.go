package handler

import (
	"auth-service/pkg/models"
	"auth-service/pkg/token"
	"auth-service/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"

	_ "auth-service/api/docs"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	AddAdmin(c *gin.Context)
	GetRole(c *gin.Context)
}

func NewAuthHandler(log *slog.Logger, authService service.AuthService) AuthHandler {
	return &authHandler{log, authService}
}

type authHandler struct {
	log *slog.Logger
	service.AuthService
}

// @Summary AddAdmin Users
// @Description create users
// @Tags Auth
// @Accept json
// @Produce json
// @Param password body models.AddingAdmin true "addid"
// @Success 200 {object} models.Success
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /admin [post]
func (h *authHandler) AddAdmin(c *gin.Context) {
	var password models.AddingAdmin

	if err := c.ShouldBindJSON(&password); err != nil {
		h.log.Error("Error in ShouldBindJSON", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{err.Error()})
		return
	}

	err := h.AuthService.AddAdmin(password)
	if err != nil {
		h.log.Error("Error in Adding admin", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Success{Message: "Admin added"})
}

// @Summary Register Users
// @Description create users
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body models.User true "register user"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /register [post]
func (h *authHandler) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error("Error in ShouldBindJSON", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{err.Error()})
		return
	}

	res, err := h.AuthService.Register(user)
	if err != nil {
		h.log.Error("Error in AuthService.Register", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Login Users
// @Description sign in user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Login body models.LoginRequest true "register user"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /login [post]
func (h *authHandler) Login(c *gin.Context) {
	var user models.LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error("Error in ShouldBindJSON", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{err.Error()})
		return
	}

	res, err := h.AuthService.Login(user)
	if err != nil {
		h.log.Error("Error in AuthService.Login", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	c.SetCookie("access_token", res.AccessToken, 3600, "", "", false, true)
	c.SetCookie("refresh_token", res.RefreshToken, 3600, "", "", false, true)

	c.JSON(http.StatusOK, res)
}

// @Summary Refreshtoken Users
// @Description sign in user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /refresh [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	value, check := c.Get("claims")
	if !check {
		h.log.Error("Error in GetClaims")
		c.JSON(http.StatusInternalServerError, models.Error{"Error in GetClaims"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		h.log.Error("Error in Claims")
		c.JSON(http.StatusInternalServerError, models.Error{"Error in Claims"})
		return
	}

	token1, err := token.GenerateAccessToken(claims.Id, claims.Role)
	if err != nil {
		h.log.Error("Error in GenerateAccessToken", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{"Error in GenerateAccessToken"})
		return
	}

	c.SetCookie("access_token", token1, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"access_token": token1})
}

// @Summary GetRole Users
// @Description get role user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /get-role [get]
func (h *authHandler) GetRole(c *gin.Context) {
	toke, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{"Error in GetRole"})
		return
	}

	claim, err := token.ExtractClaimsAccess(toke)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{"Error in ExtractClaimsAccess"})
		return
	}

	c.JSON(http.StatusOK, claim.Role)
}
