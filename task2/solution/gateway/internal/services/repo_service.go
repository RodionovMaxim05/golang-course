package services

import "gateway/internal/domain"

type RepoClient interface {
	GetRepo(url string) (domain.Repository, error)
}

type RepoService struct {
	repoClient RepoClient
}

func NewRepoService(client RepoClient) *RepoService {
	return &RepoService{repoClient: client}
}

func (rs *RepoService) GetRepo(url string) (domain.Repository, error) {
	return rs.repoClient.GetRepo(url)
}
