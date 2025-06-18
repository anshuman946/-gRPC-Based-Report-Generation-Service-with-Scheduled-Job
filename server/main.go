package main

import (
	"context"
	"fmt"
	"log"
	"net"
  "os"
  "os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	pb "report_service_updated/proto"
)

type reportServer struct {
	pb.UnimplementedReportServiceServer
	reports map[string]string
	mu      sync.RWMutex
	started time.Time
}

func (s *reportServer) GenerateReport(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.UserId == "bob" && time.Now().Unix()%3 == 0 {
		err := fmt.Errorf("temporary failure")
		log.Printf("[GenerateReport] user=%s error=%v", req.UserId, err)
		return nil, fmt.Errorf("report gen failed for %s: %w", req.UserId, err)
	}

	id := fmt.Sprintf("%s-%d", req.UserId, time.Now().UnixNano())
	s.reports[id] = fmt.Sprintf("Detailed report for %s at %s", req.UserId, time.Now().Format(time.RFC822))
	log.Printf("[GenerateReport] user=%s report_id=%s", req.UserId, id)

	return &pb.ReportResponse{
		ReportId:  id,
		StatusMsg: "Report successfully created",
		Details:   "Data snapshot taken at server timestamp",
	}, nil
}

func (s *reportServer) HealthCheck(ctx context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	uptime := time.Since(s.started).String()

	// Placeholder for RLock pattern (not accessing reports yet)
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("[HealthCheck] status=SERVING uptime=%s", uptime)
	return &pb.HealthCheckResponse{Status: "SERVING", Uptime: uptime}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("[Startup] TCP listener error: %v", err)
	}

	srv := grpc.NewServer()
	handler := &reportServer{
		reports: make(map[string]string),
		started: time.Now(),
	}

	pb.RegisterReportServiceServer(srv, handler)

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("[Startup] gRPC server running on :50051")
		go startCronJob(handler)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("[Runtime] Server crashed: %v", err)
		}
	}()

	<-stop
	log.Println("[Shutdown]  stopping server...")
	srv.GracefulStop()
	log.Println("[Shutdown] Server stopped cleanly(yeahhhh!!!)")
}
