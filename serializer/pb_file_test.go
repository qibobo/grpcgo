package serializer_test

import (
	"testing"

	"github.com/qibobo/grpcgo/models"
	"github.com/qibobo/grpcgo/serializer"
	"github.com/stretchr/testify/require"
)

func TestPB2File(t *testing.T) {
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
	err := serializer.PB2File(&p, "../tmp/pb.bin")
	require.NoError(t, err)
}
func TestFile2PB(t *testing.T) {
	p := &models.Person{}
	err := serializer.File2PB("../tmp/pb.bin", p)
	require.NoError(t, err)
}
