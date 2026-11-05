package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"ngnix-superopt/pkg/orchestrator"
)

func main() {
	target := flag.String("target", "nginx", "Target binary to optimize (nginx, envoy, or all)")
	all := flag.Bool("all", false, "Optimize all targets")
	outputJSON := flag.String("output", "", "Optional JSON file path to write report to")
	flag.Parse()

	targets := []string{*target}
	if *all || *target == "all" {
		targets = []string{"nginx", "envoy"}
	}

	for _, t := range targets {
		orch := orchestrator.NewOrchestrator(t)
		report, err := orch.RunFullPipeline()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Pipeline error for %s: %v\n", t, err)
			os.Exit(1)
		}

		bytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(bytes))

		outPath := *outputJSON
		if outPath == "" {
			outPath = fmt.Sprintf("run_pipeline_%s.json", t)
		}
		if err := os.WriteFile(outPath, bytes, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", outPath, err)
		} else {
			fmt.Printf("✅ Saved pipeline report to %s\n", outPath)
		}
	}
}
