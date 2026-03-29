package repositories

import (
	"errors"
	"go-gin-api/database"
	"go-gin-api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// Create создает нового пользователя
func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	return r.db.Create(user).Error
}

// GetAll возвращает всех пользователей с пагинацией
func (r *UserRepository) GetAll(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	
	// Получаем общее количество
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Получаем пользователей с пагинацией
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	
	return users, total, err
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail возвращает пользователя по email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(id string, updates map[string]interface{}) (*models.User, error) {
	updates["updated_at"] = time.Now()
	
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	
	if result.RowsAffected == 0 {
		return nil, nil
	}
	
	return r.GetByID(id)
}

// Delete мягкое удаление пользователя
func (r *UserRepository) Delete(id string) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	
	return nil
}

// HardDelete полное удаление пользователя
func (r *UserRepository) HardDelete(id string) error {
	result := r.db.Unscoped().Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	
	return nil
}

// GetDeleted возвращает удаленных пользователей
func (r *UserRepository) GetDeleted() ([]models.User, error) {
	var users []models.User
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users).Error
	return users, err
}

// Restore восстанавливает удаленного пользователя
func (r *UserRepository) Restore(id string) error {
	result := r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("user not found or already restored")
	}
	
	return nil
}