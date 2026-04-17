package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type SubscriptionsInfo struct {
	subscriber SubscriptionsGetter
	repoGetter RepoGetter
}

func NewSubscriptionsInfo(subscriber SubscriptionsGetter, repoGetter RepoGetter) *SubscriptionsInfo {
	return &SubscriptionsInfo{
		subscriber: subscriber,
		repoGetter: repoGetter,
	}
}

func (s *SubscriptionsInfo) Execute(ctx context.Context) ([]domain.Repository, error) {
	subscriptions, err := s.subscriber.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]domain.Repository, 0, len(subscriptions))
	for _, sub := range subscriptions {
		repoInfo, err := s.repoGetter.GetRepo(ctx, sub.Owner, sub.Repo)
		if err != nil {
			return nil, err
		}
		results = append(results, repoInfo)
	}

	return results, nil
}
