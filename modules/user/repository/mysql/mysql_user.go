package mysql

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/rachadiannovansyah/go-echo-clean-arch/domain"
	"github.com/rachadiannovansyah/go-echo-clean-arch/modules/user/repository"
	errHandle "github.com/rachadiannovansyah/go-echo-clean-arch/utils"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewMysqlUserRepository will create an object that represent the User.Repository interface
func NewMysqlUserRepository(Conn *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn}
}

func (m *mysqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.User, 0)
	for rows.Next() {
		user := domain.User{}
		err = rows.Scan(
			&user.ID,
			&user.Fullname,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.UpdatedAt,
			&user.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func (m *mysqlUserRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	query := `SELECT id, fullname, username, email, password, updated_at, created_at
  						FROM user WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)

	if err != nil && cursor != "" {
		return nil, "", errHandle.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}
