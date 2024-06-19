package storages

import "github.com/google/uuid"

type MockUserStorage struct{}

func (s *MockUserStorage) GetUserEmail(userID uuid.UUID) (string, error) {
	return "user@gmail.com", nil
}

func NewFakeUserStorage() *MockUserStorage {
	return &MockUserStorage{}
}
