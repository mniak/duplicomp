package samples

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestClientConnectingToServer(t *testing.T) {
	fakePort := gofakeit.IntRange(60000, 65535)

	go func() {
		RunServer(fakePort)
	}()

	time.Sleep(50 * time.Millisecond)

	err := RunSendPing(fakePort)
	require.NoError(t, err)
}
