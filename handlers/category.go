package handlers

import (
	"fmt"
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetCategoryHandlerData struct {
	Category models.Category  `json:"category"`
	Project  models.Project   `json:"project"`
	Trackers []models.Tracker `json:"trackers"`
	Users    []models.User    `json:"users"`
	Issues   []models.Issue   `json:"issues"`
}

// @Summary: GetCategoryHandler
// @Description: Get a category by ID
// @Tags: category
// @Produce: json
// @Param id path int true "Category ID"
// @Success 200 {object} GetCategoryHandlerData
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/{id} [get]
// @Security BearerAuth
func GetCategoryHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		category, err := models.GetCategoryByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		project, err := models.GetProjectByID(db, category.ProjectID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetCategoryHandlerData{
			Category: *category,
			Project:  *project,
			Trackers: trackers,
		}

		issues, err := models.GetIssuesByCategoryID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) > 0 {
			data.Issues = issues
		}

		users, err := models.GetUsersByCategoryID(db, category.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(users) > 0 {
			data.Users = users
		}

		c.JSON(http.StatusOK, data)
	}
}

// @Summary: CreateCategoryHandler
// @Description: Create a new category
// @Tags: category
// @Accept: json
// @Produce: json
// @Param category body models.Category true "Category"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category [post]
// @Security BearerAuth
func CreateCategoryHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		err := c.BindJSON(&category)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		id, err := models.CreateCategory(db, &category)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		category.ID = id
		c.JSON(http.StatusOK, category)
	}
}

// @Summary: UpdateCategoryHandler
// @Description: Update an existing category
// @Tags: category
// @Accept: json
// @Produce: json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/{id} [put]
// @Security BearerAuth
func UpdateCategoryHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var category models.Category
		err = c.BindJSON(&category)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != category.ID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("ID in body %d and URL %d do not match", category.ID, id)})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		err = models.UpdateCategory(db, &category)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

// @Summary: DeleteCategoryHandler
// @Description: Delete a category by ID
// @Tags: category
// @Produce: json
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/{id} [delete]
// @Security BearerAuth
func DeleteCategoryHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		err = models.DeleteCategory(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
