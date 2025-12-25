package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/pkg/apperror"
	"github.com/your-team/taskmanager-chat/backend/pkg/logging"
)

type UserService interface {
	UserRegister(users domain.User) (domain.User, error)
	UserLogin(users domain.User) (domain.TokenResponse, domain.TwoFaCodes, error)
	UserRefresh(token string) (domain.TokenResponse, error)
	UserSendEmailCode(tempToken string) error
	VerifyCode(code domain.Code) (domain.TokenResponse, error)
	EnableTwoFA(userID int64) error
	DisableTwoFA(userID int64, passqord string) error
	GetUserByID(userID int64) (domain.User, error)
}

type UsersHandler struct {
	service UserService
	logger  *logging.Logger
}

func NewUsersHandler(s UserService, l *logging.Logger) *UsersHandler {
	return &UsersHandler{
		service: s,
		logger: l,
	}
}

func (h *UsersHandler) RegisterRoutes(router *gin.RouterGroup, jwtSecret string) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/send-code", h.sendEmailToken)
		auth.POST("/verify-code", h.verifyCode)
		auth.POST("/enable-2fa", h.enableTwoFA)
		auth.POST("/disable-2fa", h.disableTwoFA)
	}
}

func (h *UsersHandler) RegisterProfileRoute(router *gin.RouterGroup) {
	router.GET("/profile", h.getProfile)
}

func (h *UsersHandler) signUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	createdUser, err := h.service.UserRegister(user)
	if err != nil {
		h.logger.Error("Failed to register user: " + err.Error())
		appErr, ok := err.(*apperror.AppError)
		if ok {
			c.JSON(http.StatusBadRequest, appErr)
		} else {
			errorMsg := "Ошибка регистрации"
			if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "пользователь с таким") {
				if strings.Contains(err.Error(), "email") {
					errorMsg = "Пользователь с таким email уже существует"
				} else if strings.Contains(err.Error(), "username") || strings.Contains(err.Error(), "именем") {
					errorMsg = "Пользователь с таким именем уже существует"
				} else {
					errorMsg = err.Error()
				}
			} else if strings.Contains(err.Error(), "all fields are required") {
				errorMsg = "Все поля обязательны для заполнения"
			} else if strings.Contains(err.Error(), "password must") {
				errorMsg = err.Error()
			} else if strings.Contains(err.Error(), "must contain letters, digits and special characters") {
				errorMsg = "Пароль должен содержать буквы, цифры и специальные символы"
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMsg,
			})
		}
		return
	}

	loginUser := domain.User{
		Email:    user.Email,
		Password: user.Password,
	}

	accessToken, tempToken, err := h.service.UserLogin(loginUser)
	if err != nil {
		h.logger.Error("Failed to login user after registration: " + err.Error())
		c.JSON(http.StatusCreated, createdUser)
		return
	}

	if tempToken.RequiresTwoFa {
		c.JSON(http.StatusOK, tempToken)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (h *UsersHandler) signIn(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат запроса",
		})
		return
	}

	accessToken, tempToken, err := h.service.UserLogin(user)
	if err != nil {
		h.logger.Error("Failed to login user: " + err.Error())
		appErr, ok := err.(*apperror.AppError)
		if ok {
			c.JSON(http.StatusUnauthorized, appErr)
		} else {
			errorMsg := "Ошибка аутентификации"
			if strings.Contains(err.Error(), "invalid credentials") {
				errorMsg = "Неверный email или пароль"
			} else if strings.Contains(err.Error(), "account is blocked") {
				errorMsg = err.Error()
			} else if strings.Contains(err.Error(), "Too many failed attempts") {
				errorMsg = "Слишком много неудачных попыток входа, аккаунт заблокирован"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errorMsg,
			})
		}
		return
	}

	if tempToken.RequiresTwoFa {
		c.JSON(http.StatusOK, tempToken)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (h *UsersHandler) refresh(c *gin.Context) {
	var req struct{ RefreshToken string `json:"refresh_token"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	token, err := h.service.UserRefresh(req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh user: " + err.Error())
		appErr, ok := err.(*apperror.AppError)
		if ok {
			c.JSON(http.StatusUnauthorized, appErr)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to refresh token",
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, token)
}

func (h *UsersHandler) sendEmailToken(c *gin.Context) {
	var req struct{ TempToken string `json:"temp_token"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	err := h.service.UserSendEmailCode(req.TempToken)
	if err != nil {
		h.logger.Error("Failed to send verification code: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to send verification code",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *UsersHandler) verifyCode(c *gin.Context) {
	var code domain.Code
	if err := c.ShouldBindJSON(&code); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	tokenRes, err := h.service.VerifyCode(code)
	if err != nil {
		h.logger.Error("Failed to verify code: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "Verification failed",
		})
		return
	}
	
	c.JSON(http.StatusOK, tokenRes)
}

func (h *UsersHandler) enableTwoFA(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	
	err := h.service.EnableTwoFA(userID.(int64))
	if err != nil {
		h.logger.Error("Failed to enable 2FA: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to enable 2FA",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *UsersHandler) disableTwoFA(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	
	var req domain.TwoFaToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
		})
		return
	}
	
	err := h.service.DisableTwoFA(userID.(int64), req.Password)
	if err != nil {
		h.logger.Error("Failed to enable 2FA: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to enable 2FA",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *UsersHandler) getProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userIDInt64 := userID.(int64)

	user, err := h.service.GetUserByID(userIDInt64)
	if err != nil {
		h.logger.Error("Failed to get user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"firstname": user.Firstname,
		"lastname": user.Lastname,
	})
}