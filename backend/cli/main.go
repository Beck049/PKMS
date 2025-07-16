package main

import (
	"fmt"
	"os"

	"pkms/backend/cli/commands"
	"pkms/backend/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Parse command
	command := os.Args[1]

	switch command {
	case "migrate":
		commands.Migrate(cfg)
	case "restore":
		commands.Restore(cfg)
	case "fix":
		commands.Fix(cfg)
	case "backup":
		commands.Backup(cfg)
	case "status":
		commands.Status(cfg)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("PKMS Database CLI Tool")
	fmt.Println("")
	fmt.Println("Usage: go run cli/main.go <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  migrate   - Run database migrations")
	fmt.Println("  restore   - Restore database from backup")
	fmt.Println("  fix       - Fix common database issues")
	fmt.Println("  backup    - Create database backup")
	fmt.Println("  status    - Check database status")
	fmt.Println("  help      - Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cli/main.go migrate")
	fmt.Println("  go run cli/main.go migrate --init=init")
	fmt.Println("  go run cli/main.go migrate --init=update")
	fmt.Println("  go run cli/main.go backup --output=backup_$(date +%Y%m%d).sql")
	fmt.Println("  go run cli/main.go restore")
	fmt.Println("  go run cli/main.go restore --force")
	fmt.Println("  go run cli/main.go restore --migrate-file=update")
	fmt.Println("  go run cli/main.go restore --force --migrate-file=update")
	fmt.Println("  go run cli/main.go fix --check-only")
}
