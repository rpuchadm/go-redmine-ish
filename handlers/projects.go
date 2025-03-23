package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProjectsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		projects, err := models.GetAllProjects(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountProjects(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"projects": projects,
			"count":    count,
		}

		c.JSON(http.StatusOK, data)
	}
}

func GetProjectHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project, err := models.GetProjectByID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"project": project,
		}

		categories, err := models.GetCategoriesByProjectID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categories) > 0 {
			data["categories"] = categories
		}

		users, err := models.GetUsersByProjectID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(users) > 0 {
			data["users"] = users
		}

		categorynumberofissues, err := models.CountIssuesByCategoryWhereProject(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categorynumberofissues) > 0 {
			data["categorynumberofissues"] = categorynumberofissues
		}

		c.JSON(http.StatusOK, data)
	}
}

func CreateProjectHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var project models.Project
		err := c.BindJSON(&project)
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

		id, err := models.CreateProject(db, &project)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		project.ID = id

		c.JSON(http.StatusCreated, project)
	}
}

func UpdateProjectHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var project models.Project
		if err := c.ShouldBindJSON(&project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != project.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID in body and URL do not match"})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		err = models.UpdateProject(db, &project)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated, err := models.GetProjectByID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

func DeleteProjectHandler(cfg *config.Config) gin.HandlerFunc {
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

		err = models.DeleteProject(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
