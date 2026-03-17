package usecases

import "collector/internal/domain"

type RepoService interface {
	GetRepo(owner, name string) (domain.Repository, error)
}

type RepoUsecases struct {
	client RepoService
}

func NewRepoUsecase(client RepoService) *RepoUsecases {
	return &RepoUsecases{client: client}
}

func (ru *RepoUsecases) Execute(owner, name string) (domain.Repository, error) {
	return ru.client.GetRepo(owner, name)
}
