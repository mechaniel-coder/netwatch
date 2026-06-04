package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mechaniel-coder/netwatch/agent/internal/collector"
	"github.com/mechaniel-coder/netwatch/agent/internal/reporter"
)

func main() {
	serverURL := flag.String("server", "", "NetWatch server URL (required)")
	token := flag.String("token", "", "Agent auth token (required)")
	interval := flag.Duration("interval", 15*time.Second, "Metric reporting interval")
	flag.Parse()

	if *serverURL == "" {
		*serverURL = os.Getenv("NETWATCH_SERVER")
	}
	if *token == "" {
		*token = os.Getenv("NETWATCH_TOKEN")
	}
	if *serverURL == "" || *token == "" {
		log.Fatal("--server and --token are required (or NETWATCH_SERVER / NETWATCH_TOKEN env vars)")
	}

	r := reporter.New(*serverURL, *token)
	agentID, err := r.Register()
	if err != nil {
		log.Fatalf("registration failed: %v", err)
	}
	log.Printf("agent registered — id: %s, reporting every %s", agentID, *interval)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-quit:
			log.Println("agent shutting down")
			return

		case <-ticker.C:
			snap, err := collector.Collect()
			if err != nil {
				log.Printf("collect error: %v", err)
				continue
			}
			if err := r.Send(agentID, snap); err != nil {
				log.Printf("send error: %v", err)
			}
		}
	}
}
