package services

import (
	"fmt"
	"main/logs"
	databaseDistributor "main/pkg/database/distributor"
	database "main/pkg/database/location"
	databaseLocation "main/pkg/database/location"
	"main/pkg/models"
	repository "main/pkg/repository/distributor"
	"strings"
)

type DistributorService struct {
	DistributorRepository repository.DistributorRepository
}

func NewDisributorService(distributorRepository repository.DistributorRepository) *DistributorService {

	newDisributorService := DistributorService{DistributorRepository: distributorRepository}

	return &newDisributorService
}

func (distributorService *DistributorService) AddDistibutorService(distributor models.Distributor) error {

	if distributor.Id == "" {
		return fmt.Errorf("distributor id required")
	}

	if len(distributor.IncludeCode) > 0 {
		for _, eachCode := range distributor.IncludeCode {

			isCorrect := checkAreaCode(eachCode)

			if !isCorrect {
				return fmt.Errorf("areaCode is not correct or is not in correct format[%v]", eachCode)
			}
		}

	}

	if len(distributor.ExcludeCode) > 0 {
		for _, eachCode := range distributor.ExcludeCode {

			isCorrect := checkAreaCode(eachCode)

			if !isCorrect {
				return fmt.Errorf("areaCode is not correct or is not in correct format[%v]", eachCode)
			}
		}

	}

	// check if  it has sub distributor, if so then it should not be its head distributor i.e sub's sub should not be cuurent ID

	if distributor.SubDistributorId != "" {

		subDistributor, error := distributorService.DistributorRepository.GetDistributor(distributor.SubDistributorId)

		if error != nil && error.Error() == "id is invalid" {
			return fmt.Errorf("sub distributor Id is invalid")
		}

		if subDistributor.SubDistributorId != "" && subDistributor.SubDistributorId == distributor.Id {
			return fmt.Errorf("deadlock")
		}

		subDistributor.HeadDistributorId = distributor.Id

		tmpDistributor, distributorError := distributorService.DistributorRepository.GetDistributor(distributor.Id)
		if distributorError != nil {
			return fmt.Errorf("distributor id is invalid")
		}
		logs.InfoLog("distributor is present on datastore.need to only assign subdistributor or update the values")
		tmpDistributor.SubDistributorId = distributor.SubDistributorId
		distributor = *tmpDistributor

		for _, areaCode := range subDistributor.IncludeCode {
			if isPresentInExclude(areaCode, &distributor) {
				return fmt.Errorf("sub distributor is not allowed to distribute the film")
			}
		}

	}

	distributor.Id = strings.TrimSpace(distributor.Id)
	distributorService.DistributorRepository.AddDistibutor(distributor)
	return nil
}

func (distributorService *DistributorService) GetDistributorById(id string) (*models.Distributor, error) {

	for _, distributor := range databaseDistributor.Distibutors {
		if distributor.Id == id {
			return distributor, nil
		}
	}
	return nil, fmt.Errorf("distributor not found with id %v", id)
}

func (distributorService *DistributorService) GetAllDistibutorService() []*models.Distributor {

	return databaseDistributor.Distibutors

}

func (distributorService *DistributorService) DistributionCheckService(request models.DistributionCheckRequest) (bool, error) {

	// if id is wrong
	distributor, err := distributorService.GetDistributorById(request.Id)
	if err != nil || distributor == nil {
		return false, fmt.Errorf("distributor id is invalid")
	}
	// check area code
	isCorrect := checkAreaCode(request.Area)
	if !isCorrect {
		return false, fmt.Errorf("area code is invalid")
	}

	lengthOfAreaCode := len(strings.Split(request.Area, "-"))
	AreaCodeArray := strings.Split(request.Area, "-")
	if lengthOfAreaCode < 1 || lengthOfAreaCode > 3 {
		return false, fmt.Errorf("areaCode is to be given in cityCode-provinceCode-countryCode format")
	}

	//only country code is provided
	if lengthOfAreaCode == 1 {
		countryCode := request.Area
		// countryAllowed := false

		currentDistributor := distributor
		if isPresentInExclude(countryCode, currentDistributor) {
			return false, nil
		}
		if !isPresentInInclude(countryCode, currentDistributor) {
			return false, nil
		}

		for currentDistributor.HeadDistributorId != "" {
			if isPresentInExclude(countryCode, currentDistributor) {
				return false, nil
			}
			currentDistributor, _ = distributorService.GetDistributorById(currentDistributor.HeadDistributorId)
		}

	} else if lengthOfAreaCode == 2 { //province and country code is provided
		provinceCode := AreaCodeArray[0]
		countryCode := AreaCodeArray[1]

		logs.InfoLog("Recieved two code area format.provinceCode [%v], countryCOde [%v]", provinceCode, countryCode)
		if !isPresentInInclude(countryCode, distributor) && !isPresentInInclude(request.Area, distributor) {
			return false, nil
		}

		if isPresentInExclude(request.Area, distributor) || isPresentInExclude(countryCode, distributor) {
			return false, nil
		}

		currentDistributor := distributor
		for currentDistributor.HeadDistributorId != "" {
			if isPresentInExclude(request.Area, distributor) || isPresentInExclude(countryCode, currentDistributor) {
				return false, nil
			}
			currentDistributor, _ = distributorService.GetDistributorById(currentDistributor.HeadDistributorId)
		}

	} else if lengthOfAreaCode == 3 { // city, province and country is provided
		cityCode := AreaCodeArray[0]
		provinceCode := AreaCodeArray[1]
		countryCode := AreaCodeArray[2]

		logs.InfoLog("Recieved three code area format.cityCode [%v], provinceCode [%v], countryCOde [%v]", cityCode, provinceCode, countryCode)

		if !isPresentInInclude(request.Area, distributor) && !isPresentInInclude(provinceCode+"-"+countryCode, distributor) && !isPresentInInclude(countryCode, distributor) {
			return false, nil
		}

		if isPresentInExclude(request.Area, distributor) || isPresentInExclude(provinceCode+"-"+countryCode, distributor) || isPresentInExclude(countryCode, distributor) {
			return false, nil
		}

		//
		currentDistributor := distributor

		for currentDistributor != nil && currentDistributor.Id != "" {

			logs.InfoLog("Checking info for ::[%v]", currentDistributor)

			if isPresentInExclude(request.Area, currentDistributor) || isPresentInExclude(provinceCode+"-"+countryCode, currentDistributor) || isPresentInExclude(countryCode, currentDistributor) {
				return false, nil
			}

			if currentDistributor.HeadDistributorId != "" {
				currentDistributor, _ = distributorService.GetDistributorById(currentDistributor.HeadDistributorId)
			} else {
				break
			}

		}
	}

	return true, nil
}

func checkAreaCode(areaCode string) bool {
	logs.InfoLog("Checking string - %v", areaCode)
	codesArray := strings.Split(areaCode, "-")
	lengthOfCode := len(codesArray)
	if lengthOfCode > 3 {
		return false
	}
	city := ""
	province := ""
	country := ""
	for index, value := range codesArray {
		if lengthOfCode == 3 && index == 0 {
			city = value
		}
		if lengthOfCode == 3 && index == 1 {
			province = value
		}
		if lengthOfCode == 3 && index == 2 {
			country = value
		}
		if lengthOfCode == 2 && index == 0 {
			province = value
		}
		if lengthOfCode == 2 && index == 1 {
			country = value
		}
		if lengthOfCode == 1 && index == 0 {
			country = value
		}
	}
	logs.InfoLog("COUNTRY:%v PROVINCE:%v, CITY:%v", country, province, city)

	//only country code is provided
	if lengthOfCode == 1 {
		if _, exists := database.AreaMap[country]; !exists {
			logs.InfoLog("Not able to get country with code %v", country)
			return false
		}
	} else if lengthOfCode == 2 {
		if _, exists := databaseLocation.AreaMap[country][province]; !exists {
			logs.InfoLog("Not able to get province with code %v", province)
			return false
		}
	} else if lengthOfCode == 3 {

		cities := databaseLocation.AreaMap[country][province]
		logs.InfoLog("Checking city [%v] in cities [%v] ", city, cities)
		found := false
		for _, value := range cities {
			if city == value {
				found = true
				break
			}
		}

		return found

	}

	return true

}

func isPresentInExclude(areaCode string, distributor *models.Distributor) bool {

	for _, eachDistributorExludeString := range distributor.ExcludeCode {
		if areaCode == eachDistributorExludeString {
			return true
		}
		// splitExcludeStringArray := strings.Split(eachDistributorExludeString, "-")

		// for _, areaCodeTarget := range splitExcludeStringArray {
		// 	if areaCode == eachDistributorExludeString {
		// 		return true
		// 	}
		// }

	}
	return false
}

func isPresentInInclude(areaCode string, distributor *models.Distributor) bool {

	for _, eachDistributorInludeString := range distributor.IncludeCode {
		if areaCode == eachDistributorInludeString {
			return true
		}
		// }
		// splitIncludeStringArray := strings.Split(eachDistributorInludeString, "-")

		// for _, areaCodeTarget := range splitIncludeStringArray {
		// 	if areaCode == areaCodeTarget {
		// 		return true
		// 	}
		// }

	}
	return false
}
