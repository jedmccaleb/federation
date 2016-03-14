package db


import (
	"database/sql"
)

type Driver interface {
	Init(url string) (err error)
	Exec(url string) (result sql.Result, err error)
	GetByStellarAddress(name, query string) (*FederationRecord, error)
	GetByAccountId(accountId, query string) (*ReverseFederationRecord, error)
}
