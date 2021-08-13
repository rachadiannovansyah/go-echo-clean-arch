package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ucase "gitlab.com/content-management-services/content-service/content/usecase"
	"gitlab.com/content-management-services/content-service/domain"
	"gitlab.com/content-management-services/content-service/domain/mocks"
)

func TestFetch(t *testing.T) {
	mockContentRepo := new(mocks.ContentRepository)
	mockContent := domain.Content{
		Title:   "Hello",
		Content: "Content",
	}

	mockListContent := make([]domain.Content, 0)
	mockListContent = append(mockListContent, mockContent)

	t.Run("success", func(t *testing.T) {
		mockContentRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListContent, "next-cursor", nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListContent))

		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockContentRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockContentRepo := new(mocks.ContentRepository)
	mockContent := domain.Content{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := domain.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockContentRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockContent, nil).Once()
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockContent.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockContentRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Content{}, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockContent.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Content{}, a)

		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockContentRepo := new(mocks.ContentRepository)
	mockContent := domain.Content{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempmockContent := mockContent
		tempmockContent.ID = 0
		mockContentRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(domain.Content{}, domain.ErrNotFound).Once()
		mockContentRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Content")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &tempmockContent)

		assert.NoError(t, err)
		assert.Equal(t, mockContent.Title, tempmockContent.Title)
		mockContentRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingContent := mockContent
		mockContentRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingContent, nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &mockContent)

		assert.Error(t, err)
		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockContentRepo := new(mocks.ContentRepository)
	mockContent := domain.Content{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockContentRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockContent, nil).Once()

		mockContentRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockContent.ID)

		assert.NoError(t, err)
		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("content-is-not-exist", func(t *testing.T) {
		mockContentRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Content{}, nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockContent.ID)

		assert.Error(t, err)
		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockContentRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Content{}, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockContent.ID)

		assert.Error(t, err)
		mockContentRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockContentRepo := new(mocks.ContentRepository)
	mockContent := domain.Content{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockContentRepo.On("Update", mock.Anything, &mockContent).Once().Return(nil)

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewContentUsecase(mockContentRepo, mockAuthorrepo, time.Second*2)

		err := u.Update(context.TODO(), &mockContent)
		assert.NoError(t, err)
		mockContentRepo.AssertExpectations(t)
	})
}
