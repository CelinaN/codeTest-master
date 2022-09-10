package service

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sports/db"
	"sports/proto/sports"
)

type Sport interface {

	// GetSports return the list of sport base on request IDs
	GetSports(ctx context.Context, in *sports.GetSportsRequest) (*sports.GetSportsResponse, error)

	// AddSport return success if added
	AddSport(ctx context.Context, in *sports.AddSportRequest) (*sports.AddSportResponse, error)
}

// sportService implements the Sport interface.
type sportService struct {
	sportRepo db.SportsRepo
}

// NewSportService instantiates and returns a new sportService.
func NewSportService(sportRepo db.SportsRepo) Sport {
	return &sportService{sportRepo}
}

// GetSports return the list of sports base on request IDs
func (s *sportService) GetSports(ctx context.Context, in *sports.GetSportsRequest) (*sports.GetSportsResponse, error) {
	sportList, err := s.sportRepo.Get(in)
	if err != nil {
		return nil, err
	}
	if sportList == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf(
			"Cannot find sport ID: %v", in.Id),
		)
	}

	return &sports.GetSportsResponse{Sport: sportList}, nil

}

// AddSport return success if added
func (s *sportService) AddSport(ctx context.Context, sportDetails *sports.AddSportRequest) (*sports.AddSportResponse, error) {
	var txt string
	txt, err := s.sportRepo.Write(sportDetails)
	if err != nil {
		txt = "Add Sport failed."
	}
	return &sports.AddSportResponse{Result: txt}, nil
}
