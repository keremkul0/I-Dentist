package userRepository_test

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/userRepository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return nil
	}
	return db
}

func cleanupDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database object")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("failed to close database connection")
	}

	err = os.Remove("test.db")
	if err != nil {
		panic("failed to remove test database")
	}
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB()
	repo := userRepository.NewUserRepository(db)

	t.Run("CreateUserRepo", func(t *testing.T) {
		user := models.User{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			Password:    "securepassword",
			CountryCode: "+1",
			PhoneNumber: "1234567890",
			NationalID:  "12345678901",
			ClinicID:    1,
		}
		createdUser, err := repo.CreateUserRepo(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		if createdUser.ID == 0 {
			t.Errorf("Expected user ID to be set, got %v", createdUser.ID)
		}
	})

	t.Run("CreateUserRepo", func(t *testing.T) {
		user := models.User{
			FirstName:   "Jane",
			LastName:    "Doe",
			Email:       "jane.doe@example.com",
			Password:    "securepassword",
			CountryCode: "+1",
			PhoneNumber: "0987654321",
			NationalID:  "10987654321",
			ClinicID:    1,
		}
		createdUser, err := repo.CreateUserRepo(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		if createdUser.ID == 0 {
			t.Errorf("Expected user ID to be set, got %v", createdUser.ID)
		}
	})

	t.Run("GetUsersRepo", func(t *testing.T) {
		users, err := repo.GetUsersRepo(1)
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %v", len(users))
		}
	})

	t.Run("GetUsersRepo", func(t *testing.T) {
		users, err := repo.GetUsersRepo(2)
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %v", len(users))
		}
	})

	t.Run("GetUserRepo", func(t *testing.T) {
		user, err := repo.GetUserRepo(1)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}
		if user.ID != 1 {
			t.Errorf("Expected user ID to be 1, got %v", user.ID)
		}
	})

	t.Run("GetUserRepo", func(t *testing.T) {
		user, err := repo.GetUserRepo(3)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if user.ID != 0 {
			t.Errorf("Expected user ID to be 0, got %v", user.ID)
		}
	})

	t.Run("GetUserByEmailRepo", func(t *testing.T) {
		user, err := repo.GetUserByEmailRepo("john.doe@example.com")
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}
		if user.Email != "john.doe@example.com" {
			t.Errorf("Expected user email to be john.doe@example.com, got %v", user.Email)
		}
	})

	t.Run("GetUserByEmailRepo", func(t *testing.T) {
		user, err := repo.GetUserByEmailRepo("nonexistent@example.com")
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if user.ID != 0 {
			t.Errorf("Expected user ID to be 0, got %v", user.ID)
		}
	})

	t.Run("UpdateUserRepo", func(t *testing.T) {
		user, err := repo.GetUserRepo(1)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}
		user.FirstName = "Johnathan"
		updatedUser, err := repo.UpdateUserRepo(user)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}
		if updatedUser.FirstName != "Johnathan" {
			t.Errorf("Expected user first name to be Johnathan, got %v", updatedUser.FirstName)
		}
	})

	t.Run("DeleteUserRepo", func(t *testing.T) {
		err := repo.DeleteUserRepo(1)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}
		users, err := repo.GetUsersRepo(1)
		if err != nil {
			t.Fatalf("Failed to get users: %v", err)
		}
		if len(users) != 1 {
			t.Errorf("Expected 1 user, got %v", len(users))
		}
	})

	t.Run("CheckUserExistRepo", func(t *testing.T) {
		user := models.UserGetModel{
			Email:       "jane.doe@example.com",
			NationalID:  "10987654321",
			PhoneNumber: "0987654321",
		}
		exists := repo.CheckUserExistRepo(user)
		if !exists {
			t.Errorf("Expected user to exist")
		}
	})
	cleanupDB(db)
}
