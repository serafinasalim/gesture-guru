package userController

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	_ "github.com/serafinasalim/gesture-guru/docs"
	"github.com/serafinasalim/gesture-guru/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func GenerateOTP(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	otp := base32.StdEncoding.EncodeToString(randomBytes)[:length]
	return otp, nil
}

func GenerateOTPToken() (string, error) {
	otp, err := GenerateOTP(6)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func SendMailOTP(body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"finasalim888@gmail.com",
		"csnk drdh zsqs gtxr",
		"smtp.gmail.com",
	)

	msg := "Subject: Verify Your GestureGuru Account\n" + body

	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"finasalim888@gmail.com",
		to,
		[]byte(msg),
	)
}

// @Summary Detail User
// @Tags Users
// @Param id path int true "userId"
// @Success 200 {object} models.User
// @Router /user/{id} [get]
func Detail(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	query := "SELECT id, username, email, profile, bio FROM users WHERE id = ?"
	err := models.DB.QueryRow(query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Profile,
		&user.Bio,
	)
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

// @Summary Register User
// @Description Sample Payload: <br> `{`<br>`"username": "serafina", `<br>` "email": "serafina@gmail.com", `<br>` "password": "123456", `<br>` "confirmPassword": "123456" `<br>`}`
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserRegister true "User Registration"
// @Router /user/register [post]
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

	// Generate OTP
	otp, err := GenerateOTPToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate OTP"})
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

	msg := "Your OTP for verification on GestureGuru is " + otp + " Please do not share this code with anyone."

	emails := []string{params.Email}

	SendMailOTP(msg, emails)

	queryInsert := "INSERT INTO users (username, email, password, otp) VALUES (?, ?, ?, ?)"
	result, err := models.DB.Exec(queryInsert, params.Username, params.Email, params.Password, otp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Get the last inserted ID
	insertId, err := result.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP has been sent. Please proceed with verification.", "user_id": insertId})
}

// @Summary Verify User
// @Description Sample Payload: <br> `{`<br>`"otp": "ABCDEF" `<br>`}`
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "userId"
// @Param new query string true "new" Enums(true,false) Default(true)
// @Param user body models.UserVerify true "User Verification"
// @Router /user/verify/{id} [put]
func Verify(c *gin.Context) {
	var params models.UserVerify

	id := c.Param("id")
	new := c.Query("new")

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

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	if err := models.DB.QueryRow(query, id).Scan(&exists); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User doesn't exist"})
		return
	}

	// Check OTP
	var correct bool
	queryCheckOTP := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND otp = ?)"
	if err := models.DB.QueryRow(queryCheckOTP, id, params.Otp).Scan(&correct); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !correct {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Wrong OTP Code"})
		return
	}

	// Update verified value in DB
	queryVerify := "UPDATE users SET verified = 1 WHERE id = ?"
	_, err := models.DB.Exec(
		queryVerify,
		id,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if new == "true" {
		// Setelah verify bisa punya personal status on lessons
		var lessonData models.Lesson
		const queryLesson = `SELECT 
						A.id,
						A.code, 
						A.title, 
						A.type, 
						A.video,
						A.duration 
					FROM lessons A`
		rows, err := models.DB.Query(queryLesson)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		defer rows.Close()

		var lessons []models.Lesson
		for rows.Next() {
			if err := rows.Scan(
				&lessonData.Id,
				&lessonData.Code,
				&lessonData.Title,
				&lessonData.Type,
				&lessonData.Video,
				&lessonData.Duration,
			); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
				return
			}
			lessons = append(lessons, lessonData)
		}

		// Prepare statement untuk bulk insert ke tabel tujuan
		stmt, err := models.DB.Prepare("INSERT INTO lesson_status (user_id, lesson_id) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for _, lesson := range lessons {
			_, err := stmt.Exec(id, lesson.Id)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Verified successfully"})
}

// @Summary Login User
// @Description Sample Payload: <br> `{`<br>`"username": "serafina", `<br>` "password": "123456" `<br>` }`
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User Login"
// @Router /user/login [post]
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

// @Summary Update User Email
// @Description Sample Payload: <br> `{`<br>` "email": "serafina@gmail.com" `<br>` }`
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "userId"
// @Param user body models.UserUpdateEmail true "User Update Email"
// @Router /user/email/{id} [put]
func UpdateEmail(c *gin.Context) {
	id := c.Param("id")
	var params models.UserUpdateEmail

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

	// Generate OTP
	otp, err := GenerateOTPToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate OTP"})
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

	msg := "Your OTP for verification on GestureGuru is " + otp + " Please do not share this code with anyone."

	emails := []string{params.Email}

	SendMailOTP(msg, emails)

	queryUpdateEmail := "UPDATE users SET email = ?, otp = ?, verified = 0 WHERE id = ?"
	result, err := models.DB.Exec(
		queryUpdateEmail,
		params.Email,
		otp,
		id,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP has been sent. Please verify to update your email."})
}

// @Summary Update User
// @Description Sample Payload: <br> `{`<br>` "profile": "profile.jpg", `<br>` "username": "serafinasalim", `<br>` "bio": "bio aq" `<br>` }`
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "userId"
// @Param user body models.UserUpdate true "User Update"
// @Router /user/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var params models.UserUpdate

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

	queryUpdate := "UPDATE users SET profile = ?, username = ?, bio = ? WHERE id = ?"
	result, err := models.DB.Exec(
		queryUpdate,
		params.Profile,
		params.Username,
		params.Bio,
		id,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User profile updated successfully"})
}

// @Summary Request OTP
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "userId"
// @Router /user/request-otp/{id} [put]
func RequestOTP(c *gin.Context) {

	userId := c.Param("id")

	// Get email from DB
	var email string
	var verified bool
	queryEmail := "SELECT email, verified FROM users WHERE id = ?"
	if err := models.DB.QueryRow(queryEmail, userId).Scan(&email, &verified); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Generate OTP
	otp, err := GenerateOTPToken()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate OTP"})
		return
	}

	msg := "Your OTP for verification on GestureGuru is " + otp + " Please do not share this code with anyone."

	emails := []string{email}

	SendMailOTP(msg, emails)

	queryUpdateOTP := "UDPATE users SET otp = ? WHERE id = ?"
	_, err = models.DB.Exec(queryUpdateOTP, otp, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP has been sent. Please proceed with verification."})
}
