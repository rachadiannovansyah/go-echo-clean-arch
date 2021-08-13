package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	contentHttp "gitlab.com/content-management-services/content-service/content/delivery/http"

	"gitlab.com/content-management-services/content-service/domain"
	"gitlab.com/content-management-services/content-service/domain/mocks"
)

func TestFetch(t *testing.T) {
	var mockContent domain.Content
	err := faker.FakeData(&mockContent)
	assert.NoError(t, err)
	mockUCase := new(mocks.ContentUsecase)
	mockListContent := make([]domain.Content, 0)
	mockListContent = append(mockListContent, mockContent)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListContent, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/contents?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := contentHttp.ContentHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchContents(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.ContentUsecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/contents?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := contentHttp.ContentHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchContents(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockContent domain.Content
	err := faker.FakeData(&mockContent)
	assert.NoError(t, err)

	mockUCase := new(mocks.ContentUsecase)

	num := int(mockContent.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockContent, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/contents/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("contents/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := contentHttp.ContentHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockContent := domain.Content{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempmockContent := mockContent
	tempmockContent.ID = 0
	mockUCase := new(mocks.ContentUsecase)

	j, err := json.Marshal(tempmockContent)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.Content")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/contents", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/contents")

	handler := contentHttp.ContentHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockContent domain.Content
	err := faker.FakeData(&mockContent)
	assert.NoError(t, err)

	mockUCase := new(mocks.ContentUsecase)

	num := int(mockContent.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/contents/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("contents/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := contentHttp.ContentHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
