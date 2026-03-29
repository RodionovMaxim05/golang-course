package usecases

import "gateway/internal/domain"

type RepoService interface {
	GetRepo(url string) (domain.Repository, error)
}

type GetRepoUsecase struct {
	client RepoService
}

func NewGetRepoUsecase(client RepoService) *GetRepoUsecase {
	return &GetRepoUsecase{client: client}
}

func (grc *GetRepoUsecase) Execute(url string) (domain.Repository, error) {
	return grc.client.GetRepo(url)
}
