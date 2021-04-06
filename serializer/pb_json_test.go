package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/serializer"
)

func TestPB2Json(t *testing.T) {
	t.Parallel()

	p := models.Person{
		Name:     "qibobo",
		Email:    "lqiyangl@gmail.com",
		IsActive: true,
		Phones: []*models.PhoneNumber{
			{
				Number: "1234567",
				Type:   models.PhoneType_HOME,
			},
			{
				Number: "13888888888",
				Type:   models.PhoneType_MOBILE,
			},
		},
	}
	err := serializer.PB2Json(&p, "../tmp/pb.json")
	require.NoError(t, err)
}
