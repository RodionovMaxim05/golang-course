package domain

// ErrorCode is a domain-level representation of why repository collection
// failed. It mirrors commonpb.ErrorCode but keeps the domain package free
// of any protobuf dependency.
type ErrorCode string

func (c ErrorCode) String() string {
	return string(c)
}

const (
	ErrorCodeUnspecified             ErrorCode = ""
	ErrorCodeRepositoryNotFound      ErrorCode = "REPOSITORY_NOT_FOUND"
	ErrorCodeGitHubRateLimitExceeded ErrorCode = "GITHUB_RATE_LIMIT_EXCEEDED"
	ErrorCodeInternalCollectorError  ErrorCode = "INTERNAL_COLLECTOR_ERROR"
)
