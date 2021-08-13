package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"gitlab.com/content-management-services/content-service/domain"
)

type contentUsecase struct {
	contentRepo    domain.ContentRepository
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

// NewContentUsecase will create new an contentUsecase object representation of domain.contentUsecase interface
func NewContentUsecase(a domain.ContentRepository, ar domain.AuthorRepository, timeout time.Duration) domain.ContentUsecase {
	return &contentUsecase{
		contentRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *contentUsecase) fillAuthorDetails(c context.Context, data []domain.Content) ([]domain.Content, error) {
	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapAuthors := map[int64]domain.Author{}

	for _, content := range data {
		mapAuthors[content.Author.ID] = domain.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *contentUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Content, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.contentRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *contentUsecase) GetByID(c context.Context, id int64) (res domain.Content, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.contentRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Content{}, err
	}
	res.Author = resAuthor
	return
}

func (a *contentUsecase) Update(c context.Context, ar *domain.Content) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.contentRepo.Update(ctx, ar)
}

func (a *contentUsecase) GetByTitle(c context.Context, title string) (res domain.Content, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.contentRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Content{}, err
	}

	res.Author = resAuthor
	return
}

func (a *contentUsecase) Store(c context.Context, m *domain.Content) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedContent, _ := a.GetByTitle(ctx, m.Title)
	if existedContent != (domain.Content{}) {
		return domain.ErrConflict
	}

	err = a.contentRepo.Store(ctx, m)
	return
}

func (a *contentUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedContent, err := a.contentRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedContent == (domain.Content{}) {
		return domain.ErrNotFound
	}
	return a.contentRepo.Delete(ctx, id)
}
