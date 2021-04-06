package service_test

import (
	"context"
	"testing"

	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/service"
	"github.com/qibobo/grpcgo/service/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestPersonServerSave(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		person *models.Person
		store  store.Store
		code   codes.Code
	}{
		{
			name: "successfully_save",
			person: &models.Person{
				Name:  "testname",
				Email: "testname@cn.com",
			},
			store: store.NewInMemoryStore(),
			code:  codes.OK,
		},
		{
			name: "failed_save",
			person: &models.Person{
				Name:  "",
				Email: "testname@cn.com",
			},
			store: store.NewInMemoryStore(),
			code:  codes.InvalidArgument,
		},
	}

	for _, tc := range testcases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			server := service.NewPersonServer(tc.store)
			resp, err := server.SavePerson(context.Background(), &models.SavePersonRequest{
				Person: tc.person,
			})
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotEmpty(t, resp.GetId())
			} else {
				require.ErrorIs(t, service.EmptyNameError, err)
				require.Empty(t, resp.GetId())
			}

		})
	}
}
