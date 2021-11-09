package repository

import (
	godis "YSGo/godis/server"
)

type Repository interface {
	Get(key string) (value string, exists bool)
	Set(key, value string)
	Flush()
}

type DefaultRepository struct {
	GodisServer *godis.GodisServer
}

var _ Repository = (*DefaultRepository)(nil)

func (repo *DefaultRepository) Get(key string) (value string, exists bool) {
	res, ok := repo.GodisServer.Get(key)

	if res == nil || !ok {
		return "", false
	}

	return res.(string), ok
}

func (repo *DefaultRepository) Set(key, value string) {
	repo.GodisServer.Set(key, value)
}

func (repo *DefaultRepository) Flush() {
	repo.GodisServer.Flush()
}
