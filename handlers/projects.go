package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetProjectsHandlerData struct {
	Projects        []models.Project `json:"projects"`
	Count           int              `json:"count"`
	IssuesNoProject []models.Issue   `json:"issues,omitempty"`
}

// @Summary: GetProjectsHandler
// @Description: Get all projects
// @Tags: projects
// @Produce: json
// @Success 200 {object} GetProjectsHandlerData
// @Failure 500 {object} map[string]string
// @Router /projects [get]
// @Security BearerAuth
func GetProjectsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			log.Println("Error GetProjectsHandler initializing database:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		projects, err := models.GetAllProjects(db)
		if err != nil {
			log.Println("Error GetProjectsHandler getting all projects:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountProjects(db)
		if err != nil {
			log.Println("Error GetProjectsHandler counting projects:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		issues, err := models.GetIssuesWhereProjectIsNull(db)
		if err != nil {
			log.Println("Error GetProjectsHandler getting issues where project is null:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetProjectsHandlerData{
			Projects:        projects,
			Count:           count,
			IssuesNoProject: issues,
		}

		c.JSON(http.StatusOK, data)
	}
}

type GetProjectHandlerData struct {
	Project                models.Project                  `json:"project"`
	Roles                  []models.Role                   `json:"roles"`
	Categories             []models.Category               `json:"categories,omitempty"`
	Users                  []models.User                   `json:"users,omitempty"`
	Members                []models.Member                 `json:"members,omitempty"`
	CategoryNumberOfIssues []models.CategoryNumberOfIssues `json:"categorynumberofissues,omitempty"`
	IssuesNoCategory       []models.Issue                  `json:"issues_no_category,omitempty"`
	Trackers               []models.Tracker                `json:"trackers"`
}

// @Summary: GetProjectHandler
// @Description: Get a project by ID
// @Tags: projects
// @Produce: json
// @Param id path int true "Project ID"
// @Success 200 {object} GetProjectHandlerData
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /project/{id} [get]
// @Security BearerAuth
func GetProjectHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			log.Println("Error initializing database:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			log.Println("Error converting project ID:", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project, err := models.GetProjectByID(db, id)
		if err != nil {
			log.Println("Error getting project by ID:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		roles, err := models.GetAllRoles(db)
		if err != nil {
			log.Println("Error getting roles:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetProjectHandlerData{
			Project:  *project,
			Roles:    roles,
			Trackers: trackers,
		}

		categories, err := models.GetCategoriesByProjectID(db, id)
		if err != nil {
			log.Println("Error getting categories by project ID:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categories) > 0 {
			data.Categories = categories
		}

		users, err := models.GetUsersByProjectID(db, id)
		if err != nil {
			log.Println("Error getting users by project ID:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(users) > 0 {
			data.Users = users
		}

		members, err := models.GetMembersByProjectID(db, id)
		if err != nil {
			log.Println("Error getting members by project ID:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(members) > 0 {
			data.Members = members
		}

		categorynumberofissues, err := models.CountIssuesByCategoryWhereProject(db, id)
		if err != nil {
			log.Println("Error counting issues by category where project:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categorynumberofissues) > 0 {
			data.CategoryNumberOfIssues = categorynumberofissues
		}

		issues, err := models.GetIssuesByProjectWhereCategoryIsNull(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) > 0 {
			data.IssuesNoCategory = issues
		}

		c.JSON(http.StatusOK, data)
	}
}

// @Summary: CreateProjectHandler
// @Description: Create a new project
// @Tags: projects
// @Accept: json
// @Produce: json
// @Param project body models.Project true "Project"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /project [post]
// @Security BearerAuth
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

// @Summary: UpdateProjectHandler
// @Description: Update a project by ID
// @Tags: projects
// @Accept: json
// @Produce: json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /project/{id} [put]
// @Security BearerAuth
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

// @Summary: UpdateProjectHandler
// @Description: Update a project by ID
// @Tags: projects
// @Accept: json
// @Produce: json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Project"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /project/{id} [put]
// @Security BearerAuth
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
