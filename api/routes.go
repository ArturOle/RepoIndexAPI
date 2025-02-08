package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "saucer_api/pkg/models"
    "saucer_api/pkg/services"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    // Health check route.
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "OK"})
    })

    // Route to trigger repository update.
    router.GET("/api/fetch-repositories", func(c *gin.Context) {
        go func() {
            // Use a service method that handles fetching and storing.
            services.FetchAndStoreRepositories()
        }()
        c.JSON(http.StatusOK, gin.H{"message": "Repository update started!"})
    })

    // Route to read repositories with query parameter.
    router.GET("/api/repositories", func(c *gin.Context) {
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
        repos, err := services.GetRepositories(limit)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, repos)
    })

    return router
}