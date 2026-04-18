package usecase

import (
	"context"

	processorpb "repo-stat/proto/processor"
)

type SubscriptionsInfoGetter interface {
	GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubsInfoRequest) (*processorpb.GetSubsInfoResponse, error)
}

type GetSubscriptionsInfo struct {
	client SubscriptionsInfoGetter
}

func NewGetSubscriptionsInfo(client SubscriptionsInfoGetter) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{client: client}
}

func (gsi *GetSubscriptionsInfo) Execute(ctx context.Context) (*processorpb.GetSubsInfoResponse, error) {
	return gsi.client.GetSubscriptionsInfo(ctx, &processorpb.GetSubsInfoRequest{})
}
