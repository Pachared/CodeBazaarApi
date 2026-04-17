package services

import (
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
)

type DownloadsService struct {
	userRepository  *repositories.UserRepository
	orderRepository *repositories.OrderRepository
}

func NewDownloadsService(
	userRepository *repositories.UserRepository,
	orderRepository *repositories.OrderRepository,
) *DownloadsService {
	return &DownloadsService{
		userRepository:  userRepository,
		orderRepository: orderRepository,
	}
}

func (s *DownloadsService) ListDownloads(currentUser *models.User) ([]contracts.DownloadLibraryItemResponse, error) {
	user, err := requireCurrentUser(currentUser)
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepository.ListDownloads(user.ID)
	if err != nil {
		return nil, err
	}

	responses := make([]contracts.DownloadLibraryItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toDownloadItemResponse(item))
	}

	return responses, nil
}

func (s *DownloadsService) MarkDownloaded(currentUser *models.User, libraryItemID string) (*contracts.MessageResponse, error) {
	user, err := requireCurrentUser(currentUser)
	if err != nil {
		return nil, err
	}

	item, err := s.orderRepository.FindDownloadForUser(user.ID, libraryItemID)
	if err != nil {
		return nil, httpx.NewAppError(404, "ไม่พบไฟล์ดาวน์โหลดที่คุณต้องการ")
	}

	now := time.Now()
	item.DownloadsCount++
	item.LastDownloadedAt = &now

	if err := s.orderRepository.SaveDownload(item); err != nil {
		return nil, err
	}

	return &contracts.MessageResponse{
		Title:       "เริ่มดาวน์โหลดแล้ว",
		Description: item.Title + " ถูกบันทึกสถานะการดาวน์โหลดเรียบร้อยแล้ว",
	}, nil
}
