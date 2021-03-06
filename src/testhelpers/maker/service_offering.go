package maker

import "cf/models"

func NewServiceOffering(label string) models.ServiceOffering {
	return models.ServiceOffering{ServiceOfferingFields: models.ServiceOfferingFields{
		Label:       label,
		Guid:        serviceOfferingGuid(),
		Description: "some service description",
	}}
}

var serviceOfferingGuid = guidGenerator("services")
