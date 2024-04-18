package limiting

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/arashi5/echo/tools/logging"
)

func TestLimiter(t *testing.T) {
	testCases := []struct {
		in  float64
		out int
	}{
		{1, 5},
		{2, 10},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		ctx = logging.WithContext(ctx, log.NewNopLogger())
		l := NewLimiter(ctx, testCase.in)

		_, listener := startMockServer(
			t,
			l.Middleware(func(ctx *fasthttp.RequestCtx) {
				ctx.Response.SetBody([]byte("hello world"))
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
			}),
		)
		defer listener.Close()

		c, err := listener.Dial()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer c.Close()

		var counter int
		for i := 0; i < 10; i++ {
			resp, err := makeGetToMockServer(t, c)

			require.NoError(t, err)
			if resp.StatusCode() == 200 {
				counter++
			}

			time.Sleep(time.Millisecond * 500)
		}

		assert.Equal(t, testCase.out, counter)
	}
}

func startMockServer(t *testing.T, handler fasthttp.RequestHandler) (*fasthttp.Server, *fasthttputil.InmemoryListener) {
	t.Helper()

	server := fasthttp.Server{Handler: handler}
	listener := fasthttputil.NewInmemoryListener()

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	return &server, listener
}

func makeGetToMockServer(t *testing.T, c net.Conn) (*fasthttp.Response, error) {
	t.Helper()

	if _, err := c.Write([]byte("GET / HTTP/1.1\r\nHost: test.dev\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	br := bufio.NewReader(c)
	var resp fasthttp.Response

	err := resp.Read(br)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return &resp, err
}
