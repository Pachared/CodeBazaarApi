package services

import "github.com/Pachared/CodeBazaarApi/internal/contracts"

import "github.com/Pachared/CodeBazaarApi/internal/repositories"

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (s *AuthService) StartGoogleAuth(intent string) (*contracts.AuthActionResponse, error) {
	user, err := s.userRepository.FindOrCreateDemoBuyer(intent)
	if err != nil {
		return nil, err
	}

	title := "เข้าสู่ระบบสำเร็จ"
	description := "คุณกำลังใช้งานบัญชีทดลองสำหรับทดสอบการเชื่อมต่อกับ API จริง"
	if intent == "register" {
		title = "สมัครสมาชิกสำเร็จ"
		description = "สร้างบัญชีผู้ใช้ทดลองเรียบร้อยแล้ว และพร้อมใช้งาน flow ฝั่งผู้ซื้อได้ทันที"
	}

	return &contracts.AuthActionResponse{
		Title:       title,
		Description: description,
		Session:     toAuthSessionUser(user),
	}, nil
}
