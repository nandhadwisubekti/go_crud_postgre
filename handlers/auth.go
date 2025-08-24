package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"go-crud-employee/database"
	"go-crud-employee/models"
	"go-crud-employee/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db         *database.DB
	jwtManager *utils.JWTManager
}

func NewAuthHandler(db *database.DB, jwtManager *utils.JWTManager) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Login authenticates user and returns JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Find user by username
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE username = $1`

	err := h.db.QueryRow(query, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Authentication failed",
				"Invalid username or password",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Authentication failed",
			"Invalid username or password",
		))
		return
	}

	// Generate JWT token
	token, expiresAt, err := h.jwtManager.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Token generation failed",
			err.Error(),
		))
		return
	}

	// Return login response
	response := models.LoginResponse{
		Token:     token,
		User:      user.ToUserInfo(),
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Login successful",
		response,
	))
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Validate password
	if !utils.IsValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid password",
			"Password must be at least 6 characters long",
		))
		return
	}

	// Check if username already exists
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", req.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.NewErrorResponse(
			"Username already exists",
			"Please choose a different username",
		))
		return
	}

	// Check if email already exists
	err = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", req.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.NewErrorResponse(
			"Email already exists",
			"Please use a different email address",
		))
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Password hashing failed",
			err.Error(),
		))
		return
	}

	// Insert new user
	var user models.User
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, username, email, created_at, updated_at`

	now := time.Now()
	err = h.db.QueryRow(query, req.Username, req.Email, hashedPassword, now, now).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Failed to create user",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"User registered successfully",
		user.ToUserInfo(),
	))
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"User not authenticated",
		))
		return
	}

	userClaims := claims.(*models.JWTClaims)

	// Get user details from database
	var user models.User
	query := `SELECT id, username, email, created_at, updated_at 
			  FROM users WHERE id = $1`

	err := h.db.QueryRow(query, userClaims.UserID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				"User not found",
				"User does not exist",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Profile retrieved successfully",
		user.ToUserInfo(),
	))
}
