package usecases

import "collector/internal/domain"

type RepoService interface {
	GetRepo(owner, name string) (domain.Repository, error)
}

type GetRepoUsecase struct {
	client RepoService
}

func NewRepoUsecase(client RepoService) *GetRepoUsecase {
	return &GetRepoUsecase{client: client}
}

func (ru *GetRepoUsecase) Execute(owner, name string) (domain.Repository, error) {
	return ru.client.GetRepo(owner, name)
}
