package user

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// UserService handles business logic related to users.
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{db: repo.db}
}

// GetUserByID fetches a user by their ID.
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Leagues").First(&user, userID).Error
	return &user, err
}

// UpdateUser updates user details in the database.
func (s *UserService) UpdateUser(userID uint, user *models.User) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(user).Error
}

// GetUserLeagues fetches all leagues that a user is in.
func (s *UserService) GetUserLeagues(userID uint) ([]models.League, error) {
	var user models.User
	err := s.db.Preload("Leagues").First(&user, userID).Error
	return user.Leagues, err
}

// GetUserPortfolios fetches all portfolios that belong to a user.
func (s *UserService) GetUserPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	err := s.db.Where("user_id = ?", userID).Find(&portfolios).Error
	return portfolios, err
}
