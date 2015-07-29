package benchmark

import (
	"os"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
	testpb "google.golang.org/grpc/benchmark/grpc_testing"
	"google.golang.org/grpc/benchmark/stats"
	"google.golang.org/grpc/grpclog"
)

func runUnary(b *testing.B, maxConcurrentCalls int) {
	s := stats.AddStats(b, 38)
	b.StopTimer()
	target, stopper := StartServer()
	defer stopper()
	conn := NewClientConn(target)
	tc := testpb.NewTestServiceClient(conn)

	// Warm up connection.
	for i := 0; i < 10; i++ {
		unaryCaller(tc)
	}
	ch := make(chan int, maxConcurrentCalls*4)
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(maxConcurrentCalls)

	// Distribute the b.N calls over maxConcurrentCalls workers.
	for i := 0; i < maxConcurrentCalls; i++ {
		go func() {
			for _ = range ch {
				start := time.Now()
				unaryCaller(tc)
				elapse := time.Since(start)
				mu.Lock()
				s.Add(elapse)
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	b.StopTimer()
	close(ch)
	wg.Wait()
	conn.Close()
}

func runStream(b *testing.B, maxConcurrentCalls int) {
	s := stats.AddStats(b, 38)
	b.StopTimer()
	target, stopper := StartServer()
	defer stopper()
	conn := NewClientConn(target)
	tc := testpb.NewTestServiceClient(conn)

	// Warm up connection.
	stream, err := tc.StreamingCall(context.Background())
	if err != nil {
		grpclog.Fatalf("%v.StreamingCall(_) = _, %v", tc, err)
	}
	for i := 0; i < 10; i++ {
		streamCaller(tc, stream)
	}

	ch := make(chan int, maxConcurrentCalls*4)
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(maxConcurrentCalls)

	// Distribute the b.N calls over maxConcurrentCalls workers.
	for i := 0; i < maxConcurrentCalls; i++ {
		go func() {
			stream, err := tc.StreamingCall(context.Background())
			if err != nil {
				grpclog.Fatalf("%v.StreamingCall(_) = _, %v", tc, err)
			}
			for _ = range ch {
				start := time.Now()
				streamCaller(tc, stream)
				elapse := time.Since(start)
				mu.Lock()
				s.Add(elapse)
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	b.StopTimer()
	close(ch)
	wg.Wait()
	conn.Close()
}
func unaryCaller(client testpb.TestServiceClient) {
	DoUnaryCall(client, 1, 1)
}

func streamCaller(client testpb.TestServiceClient, stream testpb.TestService_StreamingCallClient) {
	DoStreamingRoundTrip(client, stream, 1, 1)
}

func BenchmarkClientStreamc1(b *testing.B) {
	runStream(b, 1)
}

func BenchmarkClientStreamc8(b *testing.B) {
	runStream(b, 8)
}

func BenchmarkClientStreamc64(b *testing.B) {
	runStream(b, 64)
}

func BenchmarkClientStreamc512(b *testing.B) {
	runStream(b, 512)
}
func BenchmarkClientUnaryc1(b *testing.B) {
	runUnary(b, 1)
}

func BenchmarkClientUnaryc8(b *testing.B) {
	runUnary(b, 8)
}

func BenchmarkClientUnaryc64(b *testing.B) {
	runUnary(b, 64)
}

func BenchmarkClientUnaryc512(b *testing.B) {
	runUnary(b, 512)
}

func TestMain(m *testing.M) {
	os.Exit(stats.RunTestMain(m))
}
