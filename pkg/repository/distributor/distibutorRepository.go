package repository

import (
	"fmt"
	"main/logs"
	database "main/pkg/database/distributor"
	"main/pkg/models"
)

type DistributorRepository struct {
}

func NewDisributorRepository() *DistributorRepository {

	return &DistributorRepository{}
}

func (distributorRepository *DistributorRepository) AddDistibutor(newDisributor models.Distributor) {

	replaced := false
	for index, value := range database.Distibutors {

		if value.Id == newDisributor.Id {
			database.Distibutors[index] = &newDisributor
			replaced = true
		}
	}
	if !replaced {
		database.Distibutors = append(database.Distibutors, &newDisributor)
	}

}

func (distributorRepository *DistributorRepository) GetDistributor(id string) (*models.Distributor, error) {
	logs.InfoLog("Checking distributor with id [%v] in distribuotr list", id)
	for _, value := range database.Distibutors {
		if id == value.Id {
			return value, nil
		}
	}
	return nil, fmt.Errorf("id is invalid")
}
