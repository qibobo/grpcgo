package client

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/qibobo/grpcgo/models"
	"google.golang.org/grpc"
)

type PersonServiceClient struct {
	service models.PersonServiceClient
}

func NewPersonServiceClient(grpcClient grpc.ClientConnInterface) *PersonServiceClient {
	service := models.NewPersonServiceClient(grpcClient)
	return &PersonServiceClient{
		service: service,
	}
}

func (pc *PersonServiceClient) SavePerson(person *models.Person) (string, error) {
	resp, err := pc.service.SavePerson(context.Background(), &models.SavePersonRequest{
		Person: person,
	})
	if err != nil {
		log.Printf("failed to save person %s\n", err)
		return "", err
	}
	log.Printf("save person successfully %s", resp.GetId())
	return resp.GetId(), nil
}
func (pc *PersonServiceClient) GerPersonStream() {
	stream, err := pc.service.GetPersonStream(context.Background(), &models.GetPersonStreamRequest{})
	if err != nil {
		log.Panicf("failed to save person %s\n", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("EOF %s\n", err)
				break
			}
			log.Printf("failed to get person from stream %s\n", err)

		}
		log.Printf("get person from stream %v\n", resp)
	}
}
func (pc *PersonServiceClient) UploadImage(personId string, imageFilePath string) {
	clientStream, err := pc.service.UploadImage(context.Background())
	if err != nil {
		log.Panicf("failed to get stream for upload image %s\n", err)
		return
	}
	err = clientStream.Send(&models.UploadImageRequest{
		Data: &models.UploadImageRequest_Info{
			Info: &models.ImageInfo{
				PersonId:  personId,
				ImageType: "jpg",
			},
		},
	})
	if err != nil {
		log.Panicf("failed to send request for upload image %s\n", err)
		return
	}

	file, err := os.Open(imageFilePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()
	// imageBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Fatal("failed to read file to buffer: ", err)
	// }
	// err = clientStream.Send(&models.UploadImageRequest{
	// 	Data: &models.UploadImageRequest_ChunkData{
	// 		ChunkData: imageBytes,
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal("failed to send file: ", err)
	// }
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal("failed to read file to buffer: ", err)
		}
		err = clientStream.Send(&models.UploadImageRequest{
			Data: &models.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		})
		if err != nil {
			log.Fatal("failed to send file: ", err)
		}

	}
	uploadResp, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatal("failed to get response: ", err)
	}
	log.Printf("upload image %v\n", uploadResp)
}
