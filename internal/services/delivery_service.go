package services

import (
	"ihb-transport/internal/models"
	"ihb-transport/internal/repository"
)

type deliveryService struct {
	deliveryRepo repository.DeliveryInterface
}

func NewDeliveryService(deliveryRepo repository.DeliveryInterface) DeliveryServiceInterface {
	return &deliveryService{deliveryRepo: deliveryRepo}

}

//delivery service functions

func (s *deliveryService) CreateDeliveryRequest(delivery *models.DeliveryRequest)
