package handlers

import (
	"net/http"

	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"

	"github.com/gin-gonic/gin"
)

func InitHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		drop := true
		test := true
		sample := true

		// projects
		if drop {
			err = models.DropProjectsTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = models.CreateProjectsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if test {
			err = models.TestProjectsTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if sample {
			err = models.SampleProjects(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// trackers
		if drop {
			err = models.DropTrackersTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = models.CreateTrackersTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if test {
			err = models.TestTrackersTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if sample {
			err = models.SeedTrackers(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// roles
		if drop {
			err = models.DropRolesTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = models.CreateRolesTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if test {
			err = models.TestRolesTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if sample {
			err = models.SeedRolesTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// users
		if drop {
			err = models.DropUsersTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = models.CreateUsersTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if test {
			err = models.TestUsersTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if sample {
			err = models.SampleUsers(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// issues
		if drop {
			err = models.DropIssuesTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		err = models.CreateIssuesTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if test {
			err = models.TestIssuesTable(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if sample {
			err = models.SampleIssues(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Base de datos inicializada correctamente"})
	}
}
