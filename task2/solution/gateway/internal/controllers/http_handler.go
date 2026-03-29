package controllers

import (
	"encoding/json"
	"errors"
	"gateway/internal/domain"
	"log"
	"net/http"
)

type GetRepoUsecase interface {
	Execute(url string) (domain.Repository, error)
}

type RepoHandler struct {
	repoUsecase GetRepoUsecase
}

func NewRepoHandler(repoUsecase GetRepoUsecase) *RepoHandler {
	return &RepoHandler{repoUsecase: repoUsecase}
}

// @Summary     Get repository info
// @Description Returns repository info by GitHub URL
// @Param       url query string true "GitHub repository URL (e.g. https://github.com/golang/go)"
// @Success     200 {object} domain.Repository
// @Failure     400 {string} string "invalid request"
// @Failure     404 {string} string "repository not found"
// @Failure     500 {string} string "internal error"
// @Router      /api/v1/repo [get]
func (rh *RepoHandler) GetRepo(rw http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(rw, "url query param is required", http.StatusBadRequest)
		return
	}

	repo, err := rh.repoUsecase.Execute(url)
	if err != nil {
		mapError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(repo); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func mapError(rw http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(rw, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrRateLimited):
		http.Error(rw, err.Error(), http.StatusTooManyRequests)
	case errors.Is(err, domain.ErrInvalidArgument):
		http.Error(rw, err.Error(), http.StatusBadRequest)
	default:
		http.Error(rw, "internal server error", http.StatusInternalServerError)
	}
}
