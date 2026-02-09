package handler

import (
	"net/http"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/dto"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.authService.Signup(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or Username already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

func (h *AuthHandler) Signin(c *gin.Context) {
	var req dto.SigninRequest
	_ = c.ShouldBindJSON(&req)

	resp, err := h.authService.Signin(req)
	if err != nil {
		if err.Error() == "conflict" {
			c.JSON(http.StatusConflict, gin.H{"error": "Account logged in on another device."})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	_ = c.ShouldBindJSON(&input)

	tokens, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}
