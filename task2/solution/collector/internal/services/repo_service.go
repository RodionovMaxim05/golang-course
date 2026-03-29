package services

import "collector/internal/domain"

type RepoClient interface {
	GetRepo(owner, name string) (domain.Repository, error)
}

type RepoService struct {
	repoClient RepoClient
}

func NewRepoService(client RepoClient) *RepoService {
	return &RepoService{repoClient: client}
}

func (rs *RepoService) GetRepo(owner, name string) (domain.Repository, error) {
	return rs.repoClient.GetRepo(owner, name)
}
