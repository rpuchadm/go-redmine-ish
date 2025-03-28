package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUsersHandlerData struct {
	Users      []models.User     `json:"users"`
	Count      int               `json:"count"`
	Projects   []models.Project  `json:"projects"`
	Trackers   []models.Tracker  `json:"trackers"`
	UsersRoles []models.UserRole `json:"users_roles"`
	Roles      []models.Role     `json:"roles"`
}

// @Summary: GetUsersHandler
// @Description: Get all users
// @Tags: users
// @Produce: json
// @Success 200 {object} GetUsersHandlerData
// @Failure 500 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
func GetUsersHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		users, err := models.GetAllUsers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		roles, err := models.GetAllRoles(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user_roles, err := models.GetAllUsersRoles(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetUsersHandlerData{
			Users:      users,
			Roles:      roles,
			UsersRoles: user_roles,
		}

		c.JSON(http.StatusOK, data)
	}
}

type GetUserHandlerData struct {
	User       models.User       `json:"user"`
	Trackers   []models.Tracker  `json:"trackers"`
	Projects   []models.Project  `json:"projects,omitempty"`
	Roles      []models.Role     `json:"roles,omitempty"`
	Issues     []models.Issue    `json:"issues,omitempty"`
	Categories []models.Category `json:"categories,omitempty"`
}

// @Summary: GetUserHandler
// @Description: Get a user by ID
// @Tags: users
// @Produce: json
// @Param id path int true "User ID"
// @Success 200 {object} GetUserHandlerData
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [get]
// @Security BearerAuth
func GetUserHandler(cfg *config.Config) gin.HandlerFunc {
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

		user, err := models.GetUserByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetUserHandlerData{
			User:     *user,
			Trackers: trackers,
		}

		roles, err := models.GetRolesByUserID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(roles) > 0 {
			data.Roles = roles
		}

		issues, err := models.GetIssuesByUserID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) > 0 {
			data.Issues = issues
		}

		projects, err := models.GetProjectsByUserID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(projects) > 0 {
			data.Projects = projects
		}

		categories, err := models.GetCategoriesByUserID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(categories) > 0 {
			data.Categories = categories
		}

		c.JSON(http.StatusOK, data)
	}
}

// @Summary: CreateUserHandler
// @Description: Create a new user
// @Tags: users
// @Accept: json
// @Produce: json
// @Param user body models.User true "User"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [post]
// @Security BearerAuth
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

// @Summary: UpdateUserHandler
// @Description: Update a user by ID
// @Tags: users
// @Accept: json
// @Produce: json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [put]
// @Security BearerAuth
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

// @Summary: UpdateUserHandler
// @Description: Update a user by ID
// @Tags: users
// @Accept: json
// @Produce: json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [put]
// @Security BearerAuth
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
