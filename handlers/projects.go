package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetProjectsHandlerData struct {
	Projects []models.Project `json:"projects"`
	Count    int              `json:"count"`
}

// @Summary: GetProjectsHandler
// @Description: Get all projects
// @Tags: projects
// @Produce: json
// @Success 200 {object} GetProjectsHandlerData
// @Failure 500 {object} map[string]string
// @Router /projects [get]
// @Security Bearer
func GetProjectsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		projects, err := models.GetAllProjects(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountProjects(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetProjectsHandlerData{
			Projects: projects,
			Count:    count,
		}

		c.JSON(http.StatusOK, data)
	}
}

type GetProjectHandlerData struct {
	Project        models.Project                  `json:"project"`
	Roles          []models.Role                   `json:"roles"`
	Categories     []models.Category               `json:"categories,omitempty"`
	Users          []models.User                   `json:"users,omitempty"`
	Members        []models.Member                 `json:"members,omitempty"`
	CategoryIssues []models.CategoryNumberOfIssues `json:"categorynumberofissues,omitempty"`
}

// @Summary: GetProjectHandler
// @Description: Get a project by ID
// @Tags: projects
func GetProjectHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project, err := models.GetProjectByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		roles, err := models.GetAllRoles(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetProjectHandlerData{
			Project: *project,
			Roles:   roles,
		}

		categories, err := models.GetCategoriesByProjectID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categories) > 0 {
			data["categories"] = categories
		}

		users, err := models.GetUsersByProjectID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(users) > 0 {
			data["users"] = users
		}

		members, err := models.GetMembersByProjectID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(members) > 0 {
			data["members"] = members
		}

		categorynumberofissues, err := models.CountIssuesByCategoryWhereProject(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
