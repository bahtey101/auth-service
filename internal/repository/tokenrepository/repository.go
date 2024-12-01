package tokenrepository

import (
	"context"

	"auth-service/pkg/dbstore"
	"auth-service/tools"

	"github.com/google/uuid"
)

type TokenRepository struct {
	store dbstore.Store
}

func NewTokenRepository(store dbstore.Store) *TokenRepository {
	return &TokenRepository{
		store: store,
	}
}

func (tr *TokenRepository) Create(ctx context.Context, userID uuid.UUID, userIP, refreshToken string) error {
	const sql = `
	insert into tokens(
		user_id,
		user_ip,
		token
	) values ($1, $2, $3);
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		userID,
		userIP,
		tools.Hash(refreshToken),
	); err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) Update(ctx context.Context, userID uuid.UUID, userIP, refreshToken string) error {
	const sql = `
	update tokens
	set token = $1
	where user_id = $2
	and user_ip = $3;
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		tools.Hash(refreshToken),
		userID,
		userIP,
	); err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) Delete(ctx context.Context, userID uuid.UUID, userIP string) error {
	const sql = `
	delete from tokens
	where user_id = $1
	and user_ip = $2;
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		userID,
		userIP,
	); err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) Get(ctx context.Context, userID uuid.UUID, userIP string) (string, error) {
	var refreshToken string

	const sql = `
	select token
	from tokens
	where user_id = $1
	and user_ip = $2;
	`

	if err := tr.store.QueryRow(
		ctx,
		sql,
		userID,
		userIP,
	).Scan(
		&userID,
	); err != nil {
		return refreshToken, err
	}

	return refreshToken, nil
}
