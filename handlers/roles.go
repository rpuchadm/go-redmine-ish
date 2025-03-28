package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetRolesHandlerData struct {
	Roles []models.Role `json:"roles"`
	Count int           `json:"count"`
}

// @Summary: GetRolesHandler
// @Description: Get all roles
// @Tags: roles
// @Produce: json
// @Success 200 {object} GetRolesHandlerData
// @Failure 500 {object} map[string]string
// @Router /roles [get]
// @Security BearerAuth
func GetRolesHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		roles, err := models.GetAllRoles(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountRoles(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetRolesHandlerData{
			Roles: roles,
			Count: count,
		}

		c.JSON(http.StatusOK, data)
	}
}

// @Summary: GetRoleHandler
// @Description: Get a role by ID
// @Tags: roles
// @Produce: json
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /role/{id} [get]
// @Security BearerAuth
func GetRoleHandler(cfg *config.Config) gin.HandlerFunc {
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

		role, err := models.GetRoleByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, role)
	}
}

// @Summary: CreateRoleHandler
// @Description: Create a new role
// @Tags: roles
// @Accept: json
// @Produce: json
// @Param role body models.Role true "Role"
// @Success 201 {object} models.Role
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /role [post]
// @Security BearerAuth
func CreateRoleHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role models.Role
		if err := c.BindJSON(&role); err != nil {
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

		id, err := models.CreateRole(db, &role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		role.ID = id

		c.JSON(http.StatusCreated, role)
	}
}

func UpdateRoleHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("id")

		// pasar string id a int id
		id, err := strconv.Atoi(pid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var role models.Role
		if err := c.ShouldBindJSON(&role); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != role.ID {
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

		if err := models.UpdateRole(db, &role); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updated, err := models.GetRoleByID(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// @Summary: DeleteRoleHandler
// @Description: Delete a role by ID
// @Tags: roles
// @Produce: json
// @Param id path int true "Role ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /role/{id} [delete]
// @Security BearerAuth
func DeleteRoleHandler(cfg *config.Config) gin.HandlerFunc {
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

		err = models.DeleteRole(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
