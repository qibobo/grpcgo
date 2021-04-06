package service

import (
	"context"
	"errors"
	"log"

	"github.com/qibobo/grpcgo/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qibobo/grpcgo/service/store"
)

var (
	EmptyNameError error = errors.New("name is empty")
)

type PersonServer struct {
	Store store.Store
	models.UnimplementedPersonServiceServer
}

func NewPersonServer(store store.Store) *PersonServer {
	return &PersonServer{
		Store: store,
	}
}
func (ps *PersonServer) GetPerson(ctx context.Context, req *models.GetPersonRequest) (*models.GetPersonResponse, error) {
	id := req.GetId()
	i := ps.Store.Get(id)
	return &models.GetPersonResponse{
		Person: i.(*models.Person),
	}, nil
}
func (ps *PersonServer) SavePerson(ctx context.Context, req *models.SavePersonRequest) (*models.SavePersonResponse, error) {
	p := req.GetPerson()
	log.Printf("save person %v", p)
	if p.GetName() == "" {
		return nil, EmptyNameError
	}
	id, err := ps.Store.Save(p)
	if err != nil {
		log.Printf("failed to save person for %s\n", err)
		return nil, status.Error(codes.Internal, "failed to save person for internal error")
	}
	return &models.SavePersonResponse{
		Id: id,
	}, nil
}

func (ps *PersonServer) GetPersonStream(req *models.GetPersonStreamRequest, stream models.PersonService_GetPersonStreamServer) error {
	for _, p := range ps.Store.List() {
		stream.Send(&models.GetPersonStreamResponse{
			Person: p.(*models.Person),
		})
	}
	return nil
}
