package handlers

import "realworld-fiber-sqlc/usecase/dto/sqlc"

type HandlerBase struct {
	Queries *sqlc.Queries
}

func NewHandlerQ(queries *sqlc.Queries) *HandlerBase {
	return &HandlerBase{Queries: queries}
}
