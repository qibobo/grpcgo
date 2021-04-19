package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/service/store"
)

var (
	EmptyNameError error = errors.New("name is empty")
)

type PersonServer struct {
	Store      store.Store
	ImageStore store.ImageStore
	models.UnimplementedPersonServiceServer
}

func NewPersonServer(store store.Store, imageStore store.ImageStore) *PersonServer {
	return &PersonServer{
		Store:      store,
		ImageStore: imageStore,
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
func (ps *PersonServer) UploadImage(stream models.PersonService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("can not receive image info %s\n", err)
		return status.Errorf(codes.Unknown, "can not receive image info")
	}
	personId := req.GetInfo().GetPersonId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("upload image for person id %s and imageType is %s\n", personId, imageType)
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		ctx := stream.Context()
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Printf("upload image context timeout:%s\n", ctx.Err())
			return status.Errorf(codes.DeadlineExceeded, "upload image context timeout %s", ctx.Err())
		case context.Canceled:
			log.Printf("upload image context canceled:%s\n", ctx.Err())
			return status.Errorf(codes.Canceled, "upload image context canceled %s", ctx.Err())
		default:
		}
		req, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("can not receive image data %s\n", err)
			return status.Errorf(codes.Unknown, "can not receive data info")
		}
		receivedBytes := req.GetChunkData()
		size := len(receivedBytes)
		imageSize += size
		_, err := imageData.Write(receivedBytes)
		if err != nil {
			log.Printf("failed to write image data %s\n", err)
			return status.Errorf(codes.Unknown, "failed to write image data")
		}
	}
	_, err = ps.ImageStore.Save(personId, imageType, imageData)
	if err != nil {
		log.Printf("failed to save image %s\n", err)
		return status.Errorf(codes.Unknown, "failed to save image")
	}
	resp := models.UploadImageResponse{
		Id:   personId,
		Size: uint32(imageSize),
	}
	err = stream.SendAndClose(&resp)
	if err != nil {
		log.Printf("failed to send response %s\n", err)
		return status.Errorf(codes.Unknown, "failed to send response")
	}
	return nil
}
