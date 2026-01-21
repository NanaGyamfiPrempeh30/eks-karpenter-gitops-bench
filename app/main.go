package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	startTime = time.Now()
	version   = getEnv("APP_VERSION", "1.0.0")
	env       = getEnv("ENVIRONMENT", "development")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s | Version: %s", env, version)
	log.Printf("Architecture: %s | CPUs: %d", runtime.GOARCH, runtime.NumCPU())

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	response := map[string]interface{}{
		"message":      "EKS Karpenter GitOps Bench API",
		"version":      version,
		"environment":  env,
		"hostname":     hostname,
		"architecture": runtime.GOARCH,
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"uptime": time.Since(startTime).String(),
	})
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	response := map[string]interface{}{
		"uptime_seconds": int64(time.Since(startTime).Seconds()),
		"goroutines":     runtime.NumGoroutine(),
		"memory_mb":      mem.Alloc / 1024 / 1024,
		"architecture":   runtime.GOARCH,
		"num_cpu":        runtime.NumCPU(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}