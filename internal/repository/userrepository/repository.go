package userrepository

import (
	"context"

	"github.com/google/uuid"

	"auth-service/pkg/dbstore"
)

type UserRepository struct {
	store dbstore.Store
}

func NewUserRepository(store dbstore.Store) *UserRepository {
	return &UserRepository{
		store: store,
	}
}

func (tr *UserRepository) Create(ctx context.Context, userID uuid.UUID, email string) error {
	const sql = `
	insert into users(
		id,
		email
	) values ($1, $2);
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		userID,
		email,
	); err != nil {
		return err
	}

	return nil
}

func (tr *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var email []byte

	const sql = `
	select email
	from users
	where id = $1;
	`

	if err := tr.store.QueryRow(
		ctx,
		sql,
		userID,
	).Scan(
		&email,
	); err != nil {
		return "", err
	}

	return string(email), nil
}

func (tr *UserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	const sql = `
	delete from users
	where id = $1;
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		userID,
	); err != nil {
		return err
	}

	return nil
}
