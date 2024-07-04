package userController

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/serafinasalim/gesture-guru/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Detail(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	query := "SELECT id, username, email, password, profile FROM users WHERE id = ?"
	err := models.DB.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Profile)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
}

func Register(c *gin.Context) {
	var params models.UserRegister

	if err := c.ShouldBindJSON(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Validate the user struct
	if err := validate.Struct(&params); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": validationErrors})
		return
	}

	// Check if the email already exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	if err := models.DB.QueryRow(query, params.Email).Scan(&exists); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email already registered"})
		return
	}

	queryInsert := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := models.DB.Exec(queryInsert, params.Username, params.Email, params.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Registered successfully"})
}

func Login(c *gin.Context) {
	var user models.UserLogin

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Validate the user struct
	if err := validate.Struct(&user); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": validationErrors})
		return
	}

	// Check if the email already exists
	var storedPassword string
	query := "SELECT password FROM users WHERE username = ?"
	if err := models.DB.QueryRow(query, user.Username).Scan(&storedPassword); err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	// Compare the stored password with the provided password
	if storedPassword != user.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful"})
}
