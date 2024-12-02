package tokenrepository

import (
	"context"
	"time"

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
	where user_ip = $2;
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		tools.Hash(refreshToken),
		userIP,
	); err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) Delete(ctx context.Context, userIP string) error {
	const sql = `
	delete from tokens
	where user_ip = $1;
	`

	if _, err := tr.store.Exec(
		ctx,
		sql,
		userIP,
	); err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) GetByIP(ctx context.Context, userIP string) (uuid.UUID, string, time.Time, error) {
	var (
		userID       uuid.UUID
		refreshToken string
		createdAt    time.Time
	)

	const sql = `
	select 
		user_id,
		token,
		created_at
	from tokens
	where user_ip = $1;
	`

	if err := tr.store.QueryRow(
		ctx,
		sql,
		userIP,
	).Scan(
		&userID,
		&refreshToken,
		&createdAt,
	); err != nil {
		return userID, refreshToken, createdAt, err
	}

	return userID, refreshToken, createdAt, nil
}
