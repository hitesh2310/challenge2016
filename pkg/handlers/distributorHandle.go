package handlers

import (
	"encoding/json"
	"main/pkg/models"
	"main/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DistributorHandler struct {
	DistributorService services.DistributorService
}

func NewDistributorHandler(distributorService services.DistributorService) *DistributorHandler {
	return &DistributorHandler{DistributorService: distributorService}

}

func (handler *DistributorHandler) AddDistibutor(c *gin.Context) {
	var requestDistributor models.Distributor
	if err := c.ShouldBindJSON(&requestDistributor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.DistributorService.AddDistibutorService(requestDistributor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Distributor added successfully"})

}

func (handler *DistributorHandler) GetallDistributor(c *gin.Context) {

	allDistributor := handler.DistributorService.GetAllDistibutorService()
	resultByte, _ := json.Marshal(allDistributor)

	c.JSON(http.StatusOK, gin.H{"listOfDistributor": string(resultByte)})
}

func (handler *DistributorHandler) CheckDistributionPermission(c *gin.Context) {
	var requestDistributionCheck models.DistributionCheckRequest
	if err := c.ShouldBindJSON(&requestDistributionCheck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bool, err := handler.DistributorService.DistributionCheckService(requestDistributionCheck)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"distributionPermission": bool})

}
