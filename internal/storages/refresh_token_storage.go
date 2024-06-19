package storages

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"auth/internal/services/auth"
)

type RefreshTokenStorage struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewRefreshTokenStorage(db *sqlx.DB) *RefreshTokenStorage {
	return &RefreshTokenStorage{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *RefreshTokenStorage) Create(token *auth.RefreshToken) (uuid.UUID, error) {
	builder := s.builder.
		Insert("refresh_tokens").
		Columns(`hash, expires_at`).
		Values(token.Hash, token.ExpiresAt).
		Suffix("RETURNING \"id\"")

	query, params, err := builder.ToSql()
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "build query")
	}

	var id uuid.UUID
	err = s.db.QueryRow(query, params...).Scan(&id)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "execute query")
	}

	return id, nil
}

func (s *RefreshTokenStorage) Get(id uuid.UUID) (*auth.RefreshToken, error) {
	builder := s.builder.
		Select("id, hash, expires_at").
		From("refresh_tokens").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var refreshToken auth.RefreshToken
	err = s.db.QueryRow(query, args...).Scan(
		&refreshToken.ID, &refreshToken.Hash, &refreshToken.ExpiresAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	return &refreshToken, nil
}

func (s *RefreshTokenStorage) Delete(id uuid.UUID) error {
	builder := s.builder.
		Delete("refresh_tokens").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	err = s.db.QueryRow(query, args...).Err()
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

var _ auth.RefreshTokenStorage = &RefreshTokenStorage{}
