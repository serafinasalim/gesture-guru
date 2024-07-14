package userController

import (
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
// @Param id path int true "userId"
// @Success 200 {object} models.Lesson
// @Router /lessons/{userId} [get]
func Browse(c *gin.Context) {
	userId := c.Param("userId")

	var lessonData models.Lesson
	const queryBrowse = `SELECT 
					A.id,
					A.code, 
					A.title, 
					A.type, 
					A.video,
					A.duration 
				FROM lessons A
				JOIN lesson_status B ON A.id = B.lesson_id
				WHERE B.user_id = ?`
	rows, err := models.DB.Query(queryBrowse, userId)
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

// @Summary Save Lesson
// @Tags Users
// @Param lessonId path int true "lessonId"
// @Param userId query int true "userId"
// @Router /lessons/save/{id} [put]
func SaveLesson(c *gin.Context) {
	lessonId := c.Param("lessonId")
	userId := c.Query("userId")

	querySaveLesson := "UDPATE lesson_status SET saved = 1 WHERE user_id = ? AND lesson_id = ?"
	_, err := models.DB.Exec(querySaveLesson, userId, lessonId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Lesson saved successfully"})
}
