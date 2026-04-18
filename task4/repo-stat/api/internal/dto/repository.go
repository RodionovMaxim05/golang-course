package dto

type RepoInfoResponse struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
