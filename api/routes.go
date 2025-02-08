package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "repo_api/config"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/api/fetch-repositories", func(c *gin.Context) {
        go func() {
            StoreRepositories()
        }()
        c.JSON(http.StatusOK, gin.H{"message": "Repositories update started!"})
    })

    r.GET("/api/repositories", func(c *gin.Context) {
        limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
        reposChan := make(chan []config.Repository)
        errChan := make(chan error)

        go func() {
            repos, err := GetRepositories(limit)
            if err != nil {
                errChan <- err
                return
            }
            reposChan <- repos
        }()

        select {
        case repos := <-reposChan:
            c.JSON(http.StatusOK, repos)
        case err := <-errChan:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
    })

    return r
}