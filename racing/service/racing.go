package service

import (
	"fmt"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)

	// GetRaces return the list of races base on request IDs
	GetRaces(ctx context.Context, in *racing.GetRacesRequest) (*racing.GetRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

// ListRaces return the full list of races
func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

// GetRaces return the list of races base on request IDs
func (s *racingService) GetRaces(ctx context.Context, in *racing.GetRacesRequest) (*racing.GetRacesResponse, error) {
	races, err := s.racesRepo.Get(in)
	if err != nil {
		return nil, err
	}
	if races == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf(
			"Cannot find races ID: %v", in.Id),
		)
	}

	return &racing.GetRacesResponse{Races: races}, nil

}
