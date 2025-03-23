package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategoryHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		category, err := models.GetCategoryByID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		project, err := models.GetProjectByID(db, category.ProjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"category": category,
			"project":  project,
			"trackers": trackers,
		}

		issues, err := models.GetIssuesByCategoryID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) > 0 {
			data["issues"] = issues
		}

		users, err := models.GetUsersByCategoryID(db, category.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(users) > 0 {
			data["users"] = users
		}

		c.JSON(http.StatusOK, data)
	}
}
