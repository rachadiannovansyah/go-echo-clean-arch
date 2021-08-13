package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"gitlab.com/content-management-services/content-service/content/repository"
	contentMysqlRepo "gitlab.com/content-management-services/content-service/content/repository/mysql"
	"gitlab.com/content-management-services/content-service/domain"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockContents := []domain.Content{
		domain.Content{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.Content{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(mockContents[0].ID, mockContents[0].Title, mockContents[0].Content,
			mockContents[0].Author.ID, mockContents[0].UpdatedAt, mockContents[0].CreatedAt).
		AddRow(mockContents[1].ID, mockContents[1].Title, mockContents[1].Content,
			mockContents[1].Author.ID, mockContents[1].UpdatedAt, mockContents[1].CreatedAt)

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM contents WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := contentMysqlRepo.NewMysqlContentRepository(db)
	cursor := repository.EncodeCursor(mockContents[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM contents WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := contentMysqlRepo.NewMysqlContentRepository(db)

	num := int64(5)
	aContents, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, aContents)
}

func TestStore(t *testing.T) {
	now := time.Now()
	ar := &domain.Content{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT contents SET title=\\? , content=\\? , author_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := contentMysqlRepo.NewMysqlContentRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM contents WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := contentMysqlRepo.NewMysqlContentRepository(db)

	title := "title 1"
	aContents, err := a.GetByTitle(context.TODO(), title)
	assert.NoError(t, err)
	assert.NotNil(t, aContents)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM contents WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := contentMysqlRepo.NewMysqlContentRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	ar := &domain.Content{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE contents set title=\\?, content=\\?, author_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := contentMysqlRepo.NewMysqlContentRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
