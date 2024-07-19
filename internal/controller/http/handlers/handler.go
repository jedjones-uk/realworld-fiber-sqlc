package handlers

import (
	"realworld-fiber-sqlc/pkg/logger"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

type HandlerBase struct {
	Queries *sqlc.Queries
	*logger.Logger
}

func NewHandlerQ(queries *sqlc.Queries, l *logger.Logger) *HandlerBase {
	return &HandlerBase{
		Queries: queries,
		Logger:  l,
	}
}
