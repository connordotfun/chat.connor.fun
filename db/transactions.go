package db

import "github.com/jmoiron/sqlx"

type sqlxTransaction struct {
	Repository
	err error
	tx *sqlx.Tx
}

func (t *sqlxTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *sqlxTransaction) Rollback() error {
	return t.tx.Rollback()
}


type sqlxTransactionalRepository struct {
	RepositoryImpl
}

func NewTransactionalRepository(db *sqlx.DB, users UserRepository, roles RolesRepository,
		rooms RoomsRepository, messages MessagesRepository, verifications VerificationCodeRepository) *sqlxTransactionalRepository {
	return &sqlxTransactionalRepository{
		RepositoryImpl: RepositoryImpl{
			Source: db,
			UsersRepo: users,
			MessagesRepo: messages,
			RolesRepo: roles,
			RoomsRepo: rooms,
			VerificationsRepo: verifications,
		},
	}
}

func (tr *sqlxTransactionalRepository) CreateTransaction() Transaction {
	rootDb, ok := tr.Source.(*sqlx.DB)
	if !ok {
		panic("non-transactional source used for transactional repository!")
	}

	tx := rootDb.MustBegin()

	return &sqlxTransaction{
		Repository: tr.NewFromSource(tx),
		tx: tx,
	}
}
