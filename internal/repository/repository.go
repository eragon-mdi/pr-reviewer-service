package repository

import (
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/storage"
	sqlrepo "github.com/eragon-mdi/pr-reviewer-service/internal/repository/sql"
	"github.com/eragon-mdi/pr-reviewer-service/internal/service"
)

func New(s storage.Storage) service.Repository {
	return &repository{
		SqlRepo: sqlrepo.New(s.SQL()),
	}
}

type repository struct {
	sqlrepo.SqlRepo
}
