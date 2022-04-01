package resolvers

import "github.com/brice-74/golang-base-api/internal/api/application"

type Root struct {
	App *application.Application
}

type SortParams struct {
	Sort string
}

type PaginationParams struct {
	Offset int32
	Limit  int32
}

type ResolverParams struct {
	SortParams
	PaginationParams
}
