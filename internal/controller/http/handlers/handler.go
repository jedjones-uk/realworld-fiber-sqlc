package handlers

import (
	"realworld-fiber-sqlc/pkg/logger"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

type HandlerBase struct {
	Queries sqlc.Querier
	Logger  logger.Interface
}

func NewHandlerQ(queries sqlc.Querier, l logger.Interface) *HandlerBase {
	return &HandlerBase{
		Queries: queries,
		Logger:  l,
	}
}
