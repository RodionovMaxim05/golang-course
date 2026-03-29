package usecase

import "repo-stat/collector/internal/domain"

type RepoService interface {
	GetRepo(owner, name string) (domain.Repository, error)
}

type GetRepoUsecase struct {
	client RepoService
}

func NewRepoUsecase(client RepoService) *GetRepoUsecase {
	return &GetRepoUsecase{client: client}
}

func (gru *GetRepoUsecase) Execute(owner, name string) (domain.Repository, error) {
	return gru.client.GetRepo(owner, name)
}
