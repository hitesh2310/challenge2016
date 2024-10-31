package router

import (
	"fmt"
	"main/pkg/handlers"
	repository "main/pkg/repository/distributor"
	"main/pkg/services"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	router := gin.Default()
	addRoutes(router)
	return router

}

func addRoutes(r *gin.Engine) {
	fmt.Println("here!")
	distributorRepository := repository.NewDisributorRepository()
	distributorService := services.NewDisributorService(*distributorRepository)
	distributorHandlers := handlers.NewDistributorHandler(*distributorService)

	r.POST("/add", distributorHandlers.AddDistibutor)
	r.GET("/all", distributorHandlers.GetallDistributor)
	r.GET("/check", distributorHandlers.CheckDistributionPermission)
}
