package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsersHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		users, err := models.GetAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		roles, err := models.GetAllRoles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user_roles, err := models.GetAllUsersRoles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"users":      users,
			"roles":      roles,
			"user_roles": user_roles,
		}

		c.JSON(http.StatusOK, data)
	}
}

func GetUserHandler(cfg *config.Config) gin.HandlerFunc {
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

		user, err := models.GetUserByID(db, id)
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
			"user":     user,
			"trackers": trackers,
		}

		roles, err := models.GetRolesByUserID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(roles) > 0 {
			data["roles"] = roles
		}

		issues, err := models.GetIssuesByUserID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) > 0 {
			data["issues"] = issues
		}

		projects, err := models.GetProjectsByUserID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(projects) > 0 {
			data["projects"] = projects
		}

		categories, err := models.GetCategoriesByUserID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categories) > 0 {
			data["categories"] = categories
		}

		c.JSON(http.StatusOK, data)
	}
}

func CreateUserHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.BindJSON(&user)
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

		id, err := models.CreateUser(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.ID = id

		c.JSON(http.StatusCreated, user)
	}
}

func UpdateUserHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != user.ID {
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

		err = models.UpdateUser(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUserHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")
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

		err = models.DeleteUser(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
