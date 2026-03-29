package services

import (
	"errors"
	"go-gin-api/dto"
	"go-gin-api/models"
	"go-gin-api/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, _ := s.repo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetAllUsers возвращает всех пользователей с пагинацией
func (s *UserService) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	offset := (page - 1) * pageSize
	return s.repo.GetAll(pageSize, offset)
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser обновляет пользователя
func (s *UserService) UpdateUser(id string, req *dto.UpdateUserRequest) (*models.User, error) {
	// Проверяем существование пользователя
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	// Если обновляется email, проверяем его уникальность
	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.repo.GetByEmail(*req.Email)
		if existingUser != nil {
			return nil, errors.New("email already in use")
		}
	}
	
	// Формируем map с обновлениями
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Age != nil {
		updates["age"] = *req.Age
	}
	
	return s.repo.Update(id, updates)
}

// DeleteUser удаляет пользователя (мягкое удаление)
func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

// HardDeleteUser полностью удаляет пользователя
func (s *UserService) HardDeleteUser(id string) error {
	return s.repo.HardDelete(id)
}

// RestoreUser восстанавливает удаленного пользователя
func (s *UserService) RestoreUser(id string) error {
	return s.repo.Restore(id)
}