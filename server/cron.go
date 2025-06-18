package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	pb "report_service_updated/proto"
)

func startCronJob(s *reportServer) {
	usersEnv := os.Getenv("REPORT_USERS")
	if usersEnv == "" {
		log.Println("[CRON] No users defined in REPORT_USERS env, defaulting to test users.")
		usersEnv = "anshu,ayush,yash"
	}
	users := strings.Split(usersEnv, ",")

	c := cron.New()
	jobID, err := c.AddFunc("@every 10s", func() {
		log.Printf("[CRON] Triggered at %s", time.Now().Format(time.RFC3339))
		for _, user := range users {
			user = strings.TrimSpace(user)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			resp, err := s.GenerateReport(ctx, &pb.ReportRequest{UserId: user})
			if err != nil {
				log.Printf("[CRON] Failed to generate report for user=%s: %v", user, err)
			} else {
				log.Printf("[CRON] Report created for user=%s report_id=%s", user, resp.ReportId)
			}
		}
	})

	if err != nil {
		log.Fatalf("[CRON] Setup failed: %v", err)
	}

	log.Printf("[CRON] Job scheduled successfully (ID=%d)", jobID)
	c.Start()
}
