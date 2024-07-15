package userController

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	_ "github.com/serafinasalim/gesture-guru/docs"
	"github.com/serafinasalim/gesture-guru/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// @Summary Browse Lesson
// @Tags Lessons
// @Param lesson body models.LessonBrowse true "LessonBrowse"
// @Success 200 {object} models.Lesson
// @Router /lessons [post]
func Browse(c *gin.Context) {
	var params models.LessonBrowse

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

	var lessonData models.Lesson
	const queryBrowse = `SELECT 
							A.id,
							A.code, 
							A.title, 
							A.type, 
							A.video,
							A.duration,
							B.saved,
							B.status
						FROM lessons A
						JOIN lesson_status B ON A.id = B.lesson_id
						WHERE B.user_id = ?
						ORDER BY B.last_seen DESC, A.id ASC `
	rows, err := models.DB.Query(queryBrowse, params.UserId)
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
			&lessonData.Saved,
			&lessonData.Status,
		); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		lessons = append(lessons, lessonData)
	}

	if err := rows.Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if len(lessons) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "No lessons found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "lessons": lessons})
}

// @Summary Detail Lesson
// @Tags Lessons
// @Param lesson body models.LessonBrowse true "LessonBrowse"
// @Router /lesson [post]
func Detail(c *gin.Context) {
	var params models.LessonBrowse

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

	var lesson models.Lesson
	queryDetail := `SELECT 
						A.id,
						A.code, 
						A.title, 
						A.type, 
						A.video,
						A.duration,
						B.saved,
						B.status
					FROM lessons A
					JOIN lesson_status B ON A.id = B.lesson_id
					WHERE B.user_id = ? AND lesson_id = ?`
	err := models.DB.QueryRow(queryDetail, params.UserId, params.LessonId).Scan(
		&lesson.Id,
		&lesson.Code,
		&lesson.Title,
		&lesson.Type,
		&lesson.Video,
		&lesson.Duration,
		&lesson.Saved,
		&lesson.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	// Update kolom last_seen dengan timestamp saat ini
	queryUpdateLastSeen := "UPDATE lesson_status SET last_seen = CURRENT_TIMESTAMP WHERE user_id = ? AND lesson_id = ?"
	_, err = models.DB.Exec(queryUpdateLastSeen, params.UserId, params.LessonId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "lesson": lesson})
}

// @Summary Save Lesson
// @Tags Lessons
// @Param lessonId path int true "lessonId"
// @Param userId path int true "userId"
// @Router /lesson/save/{userId}/{lessonId} [put]
func SaveLesson(c *gin.Context) {
	lessonId := c.Param("lessonId")
	userId := c.Param("userId")

	querySaveLesson := "UPDATE lesson_status SET saved = 1 WHERE user_id = ? AND lesson_id = ?"
	_, err := models.DB.Exec(querySaveLesson, userId, lessonId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Lesson saved successfully"})
}

// @Summary Attempt Lesson
// @Tags Lessons
// @Param lessonId path int true "lessonId"
// @Param userId path int true "userId"
// @Param user body models.LessonAttempt true "Lesson Attempt"
// @Router /lesson/attempt/{userId}/{lessonId} [put]
func AttemptLesson(c *gin.Context) {
	lessonId := c.Param("lessonId")
	userId := c.Param("userId")

	var params models.LessonAttempt

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

	// Get Attempt Number
	var attemptNumber int
	queryAttemptNumber := "SELECT attempt_number FROM lesson_scores WHERE user_id = ? AND lesson_id = ? LIMIT 1"
	err := models.DB.QueryRow(queryAttemptNumber, userId, lessonId).Scan(
		&attemptNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			attemptNumber = 1
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	queryAttemptLesson := "INSERT INTO lesson_scores (user_id, lesson_id, attempt_number, score) VALUES (?, ?, ?, ?)"
	_, err = models.DB.Exec(queryAttemptLesson, userId, lessonId, attemptNumber, params.Score)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var status string
	if attemptNumber == 1 {
		if params.Score >= 75 {
			status = "Completed"
		} else {
			status = "Failed"
		}
	} else {
		// Get status
		var statusBefore string
		queryGetStatus := "SELECT status FROM lesson_status WHERE user_id = ? AND lesson_id = ?"
		err := models.DB.QueryRow(queryGetStatus, userId, lessonId).Scan(
			&statusBefore,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			}
			return
		}

		if statusBefore == "Failed" && params.Score >= 75 {
			status = "Completed"
		} else {
			status = statusBefore
		}

	}

	queryUpdateStatus := "UPDATE lesson_status SET status = ? WHERE user_id = ? AND lesson_id = ?"
	_, err = models.DB.Exec(queryUpdateStatus, status, userId, lessonId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Lesson attempted"})
}

// @Summary Attempt Lesson
// @Tags Lessons
// @Param user body models.LessonBrowse true "Lesson Browse"
// @Router /lesson/achievement [post]
// func Achievement(c *gin.Context) {
// 	var params models.LessonBrowse

// 	if err := c.ShouldBindJSON(&params); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	// Validate the user struct
// 	if err := validate.Struct(&params); err != nil {
// 		var validationErrors []string
// 		for _, err := range err.(validator.ValidationErrors) {
// 			validationErrors = append(validationErrors, err.Error())
// 		}
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": validationErrors})
// 		return
// 	}

// 	// Get Attempt Number
// 	var attemptNumber int
// 	queryAttemptNumber := "SELECT attempt_number FROM lesson_scores WHERE user_id = ? AND lesson_id = ? LIMIT 1"
// 	err := models.DB.QueryRow(queryAttemptNumber, userId, lessonId).Scan(
// 		&attemptNumber,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			attemptNumber = 1
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 			return
// 		}
// 	}

// 	queryAttemptLesson := "INSERT INTO lesson_scores (user_id, lesson_id, attempt_number, score) VALUES (?, ?, ?, ?)"
// 	_, err = models.DB.Exec(queryAttemptLesson, userId, lessonId, attemptNumber, params.Score)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	var status string
// 	if attemptNumber == 1 {
// 		if params.Score >= 75 {
// 			status = "Completed"
// 		} else {
// 			status = "Failed"
// 		}
// 	} else {
// 		// Get status
// 		var statusBefore string
// 		queryGetStatus := "SELECT status FROM lesson_status WHERE user_id = ? AND lesson_id = ?"
// 		err := models.DB.QueryRow(queryGetStatus, userId, lessonId).Scan(
// 			&statusBefore,
// 		)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
// 			} else {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 			}
// 			return
// 		}

// 		if statusBefore == "Failed" && params.Score >= 75 {
// 			status = "Completed"
// 		} else {
// 			status = statusBefore
// 		}

// 	}

// 	queryUpdateStatus := "UPDATE lesson_status SET status = ? WHERE user_id = ? AND lesson_id = ?"
// 	err = models.DB.QueryRow(queryUpdateStatus, userId, lessonId).Scan(
// 		status,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Lesson attempted"})
// }
