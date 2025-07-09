package containersuseconnect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCredentialsNotFoundError(t *testing.T) {
	err := &SuseConnectError{
		ErrorCode: CredentialsNotFoundError,
	}

	assert.True(t, IsCredentialsNotFoundError(err))
}

func TestIsNotCredentialsNotFoundError(t *testing.T) {
	err := &SuseConnectError{
		ErrorCode: InvalidCredentialsError,
	}

	assert.False(t, IsCredentialsNotFoundError(err))
}
