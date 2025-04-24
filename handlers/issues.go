package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetIssuesHandlerData struct {
	Issues []models.Issue `json:"issues"`
}

// @Summary: InitHandler
// @Description: Initialize the database
// @Tags: init
// @Produce: json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /init [get]
// @Security BearerAuth
func GetIssuesHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		issues, err := models.GetAllIssues(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetIssuesHandlerData{
			Issues: issues,
		}

		c.JSON(http.StatusOK, data)
	}
}

type GetIssueHandlerData struct {
	Issue      *models.Issue     `json:"issue,omitempty"`
	Trackers   []models.Tracker  `json:"trackers"`
	Project    *models.Project   `json:"project,omitempty"`
	Users      []models.User     `json:"users,omitempty"`
	Categories []models.Category `json:"categories,omitempty"`
	Comments   []models.Comment  `json:"comments,omitempty"`
}

// @Summary: GetIssueHandler
// @Description: Get an issue by ID
// @Tags: issues
// @Produce: json
// @Param id path int true "Issue ID"
// @Success 200 {object} GetIssueHandlerData
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issue/{id} [get]
// @Security BearerAuth
func GetIssueHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		project_id := 0
		qproject_id := c.Query("project_id")
		if qproject_id != "" {
			var err error
			project_id, err = strconv.Atoi(qproject_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

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

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetIssueHandlerData{
			Trackers: trackers,
		}

		if id > 0 {

			issue, err := models.GetIssueByID(db, id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if issue == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
				return
			}

			data.Issue = issue

			if issue.ProjectID != nil {
				project_id = *issue.ProjectID
			}
		}

		if project_id != 0 {
			project, err := models.GetProjectByID(db, project_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			data.Project = project

			categories, err := models.GetCategoriesByProjectID(db, project_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if len(categories) > 0 {
				data.Categories = categories
			}
		}

		if id > 0 {
			users, err := models.GetUsersByIssueID(db, id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if len(users) > 0 {
				data.Users = users
			}

			comments, err := models.GetCommentsByIssueID(db, id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if len(comments) > 0 {
				data.Comments = comments
			}
		}

		c.JSON(http.StatusOK, data)
	}
}

// @Summary: CreateIssueHandler
// @Description: Create a new issue
// @Tags: issues
// @Accept: json
// @Produce: json
// @Param issue body models.Issue true "Issue"
// @Success 201 {object} models.Issue
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issue [post]
// @Security BearerAuth
func CreateIssueHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var issue models.Issue
		if err := c.ShouldBindJSON(&issue); err != nil {
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

		id, err := models.CreateIssue(db, &issue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		issue.ID = id

		c.JSON(http.StatusCreated, issue)
	}
}

// @Summary: UpdateIssueHandler
// @Description: Update an issue by ID
// @Tags: issues
// @Accept: json
// @Produce: json
// @Param id path int true "Issue ID"
// @Param issue body models.Issue true "Issue"
// @Success 200 {object} models.Issue
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issue/{id} [put]
// @Security BearerAuth
func UpdateIssueHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var issue models.Issue
		if err := c.ShouldBindJSON(&issue); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != issue.ID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID in body and URL do not match"})
			return
		}

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		if err := models.UpdateIssue(db, &issue); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated, err := models.GetIssueByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// @Summary: UpdateIssueHandler
// @Description: Update an issue by ID
// @Tags: issues
// @Accept: json
// @Produce: json
// @Param id path int true "Issue ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issue/{id} [delete]
// @Security BearerAuth
func DeleteIssueHandler(cfg *config.Config) gin.HandlerFunc {
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

		if err := models.DeleteIssue(db, id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
