package samples

import (
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientConnectingToServer(t *testing.T) {
	fakePort := gofakeit.IntRange(60000, 65535)

	go func() {
		RunServer(WithPort(fakePort))
	}()

	time.Sleep(50 * time.Millisecond)

	phrase := gofakeit.SentenceSimple()
	pong, err := RunSendPing(phrase, WithPort(fakePort))

	require.NoError(t, err)
	require.NotNil(t, pong)

	assert.Equal(t, phrase, *pong.OriginalMessage)
	assert.Equal(t, strings.ToUpper(phrase), *pong.CapitalizedMessage)
}
