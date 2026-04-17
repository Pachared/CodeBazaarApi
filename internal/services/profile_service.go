package services

import "github.com/Pachared/CodeBazaarApi/internal/contracts"

import "github.com/Pachared/CodeBazaarApi/internal/models"

import "github.com/Pachared/CodeBazaarApi/internal/repositories"

type ProfileService struct {
	userRepository *repositories.UserRepository
}

func NewProfileService(userRepository *repositories.UserRepository) *ProfileService {
	return &ProfileService{userRepository: userRepository}
}

func (s *ProfileService) GetProfile(currentUser *models.User) (*contracts.AuthSessionUser, error) {
	user, err := s.userRepository.ResolveOrDefaultBuyer(currentUser)
	if err != nil {
		return nil, err
	}

	return toAuthSessionUser(user), nil
}

func (s *ProfileService) UpdateProfile(currentUser *models.User, input contracts.ProfileUpdateRequest) (*contracts.AuthSessionUser, error) {
	user, err := s.userRepository.ResolveOrDefaultBuyer(currentUser)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	user.PhoneNumber = input.PhoneNumber
	user.StoreName = input.StoreName
	user.SavedCardHolderName = input.SavedCardHolderName
	user.SavedCardNumber = input.SavedCardNumber
	user.SavedCardExpiry = input.SavedCardExpiry
	user.BankName = input.BankName
	user.BankAccountNumber = input.BankAccountNumber
	user.BankBookImageName = input.BankBookImageName
	user.BankBookImageURL = input.BankBookImageURL
	user.IdentityCardImageName = input.IdentityCardImageName
	user.IdentityCardImageURL = input.IdentityCardImageURL
	user.NotifyOrders = input.NotifyOrders
	user.NotifyMarketplace = input.NotifyMarketplace

	if err := s.userRepository.Save(user); err != nil {
		return nil, err
	}

	return toAuthSessionUser(user), nil
}
