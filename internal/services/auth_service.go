package services

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
)

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

type googleUserInfoResponse struct {
	Sub           string      `json:"sub"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	EmailVerified interface{} `json:"email_verified"`
}

type AuthService struct {
	userRepository *repositories.UserRepository
	httpClient     *http.Client
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
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

func (s *AuthService) ExchangeGoogleSession(accessToken string, intent string) (*contracts.AuthActionResponse, error) {
	request, err := http.NewRequest(http.MethodGet, googleUserInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, httpx.NewAppError(http.StatusBadGateway, "ไม่สามารถเชื่อมต่อ Google เพื่อยืนยันตัวตนได้")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, httpx.NewAppError(http.StatusUnauthorized, "Google access token ไม่ถูกต้องหรือหมดอายุแล้ว")
	}

	var userInfo googleUserInfoResponse
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, httpx.NewAppError(http.StatusBadGateway, "ไม่สามารถอ่านข้อมูลผู้ใช้จาก Google ได้")
	}

	if strings.TrimSpace(userInfo.Sub) == "" || strings.TrimSpace(userInfo.Email) == "" {
		return nil, httpx.NewAppError(http.StatusUnauthorized, "Google ไม่ได้ส่งข้อมูลบัญชีที่จำเป็นกลับมา")
	}

	emailVerified := false
	switch value := userInfo.EmailVerified.(type) {
	case bool:
		emailVerified = value
	case string:
		emailVerified = strings.EqualFold(value, "true")
	}

	if !emailVerified {
		return nil, httpx.NewAppError(http.StatusUnauthorized, "บัญชี Google นี้ยังไม่ได้ยืนยันอีเมล")
	}

	user, err := s.userRepository.FindOrCreateExternalUser(
		"google-"+strings.TrimSpace(userInfo.Sub),
		strings.TrimSpace(strings.ToLower(userInfo.Email)),
		strings.TrimSpace(userInfo.Name),
		"google",
		"buyer",
	)
	if err != nil {
		return nil, err
	}

	title := "เข้าสู่ระบบสำเร็จ"
	description := "บัญชี Google ของคุณถูกเชื่อมเข้ากับ CodeBazaar เรียบร้อยแล้ว"
	if intent == "register" {
		title = "สมัครสมาชิกสำเร็จ"
		description = "สร้างบัญชีผู้ใช้จาก Google เรียบร้อยแล้ว สามารถใช้งานต่อได้ทันที"
	}

	return &contracts.AuthActionResponse{
		Title:       title,
		Description: description,
		Session:     toAuthSessionUser(user),
	}, nil
}
