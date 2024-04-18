package echo

import (
	"context"
	"github.com/arashi5/echo/internal/repository"
	"github.com/pkg/errors"
)

type echoService struct {
	repo *repository.Instance
}

func NewEchoService(repo *repository.Instance) Service {
	return &echoService{repo: repo}
}

func (s *echoService) CreateEcho(ctx context.Context, req *CreateEchoRequest) (resp *CreateEchoResponse, err error) {
	repo, err := s.repo.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "repo.New")
	}
	response := CreateEchoResponse{}
	err = repo.Transaction(ctx, func() error {
		if response.Id, err = repo.GetEchoRepository().CreateOrUpdate(ctx, req.Data.Dto()); err != nil {
			return errors.Wrap(err, "repo.GetEchoRepository")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Transaction")
	}
	return &response, nil
}

func (s *echoService) GetEcho(ctx context.Context, req *GetEchoListRequest) (resp *GetEchoListResponse, err error) {
	repo, err := s.repo.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "repo.New")
	}
	response := GetEchoListResponse{}
	err = repo.Transaction(ctx, func() error {
		result, err := repo.GetEchoRepository().FindAll(ctx)
		if err != nil {
			return errors.Wrap(err, "repo.GetEchoRepository")
		}
		response = make(GetEchoListResponse, 0, len(result))
		for i := range result {
			response = append(response, (echoRepo)(result[i]).Dto())
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Transaction")
	}
	return &response, nil
}
