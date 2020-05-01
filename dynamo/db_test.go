package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func Test_ConnectString(t *testing.T) {
	s := "http://127.0.0.1:8000?id=123&secret=456&token=789&region=localhost"
	cfg, err := ParseConnectString(s)
	assert.NoError(t, err)
	assert.Equal(t, &aws.Config{
		Credentials: credentials.NewStaticCredentials("123", "456", "789"),
		Endpoint:    aws.String("http://127.0.0.1:8000"),
		Region:      aws.String("localhost"),
	}, cfg)
}
