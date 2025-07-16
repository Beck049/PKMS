package commands

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pkms/backend/config"

	_ "github.com/go-sql-driver/mysql"
)

// Migrate runs database migrations
func Migrate(cfg *config.Config) {
	flagSet := flag.NewFlagSet("migrate", flag.ExitOnError)
	initFlag := flagSet.String("init", "init", "SQL file name (without .sql extension) in db/ directory")
	flagSet.Parse(os.Args[2:])

	sqlFile := fmt.Sprintf("db/%s.sql", *initFlag)
	fmt.Printf("Running migration using %s...\n", sqlFile)

	// 1. 讀取 SQL 檔案
	content, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("Failed to read SQL file:", err)
	}

	// 2. 連接資料庫
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 3. 分割並執行每條 SQL
	statements := splitSQL(string(content))
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Failed to execute statement %d: %v\nSQL: %s\n", i+1, err, stmt)
		}
	}

	fmt.Println("Database migration completed successfully!")
}

// Backup creates a database backup (Under Developed)
func Backup(cfg *config.Config) {
	flagSet := flag.NewFlagSet("backup", flag.ExitOnError)
	outputFlag := flagSet.String("output", "", "Output file for backup")
	flagSet.Parse(os.Args[2:])

	if *outputFlag == "" {
		// Generate default filename with timestamp
		timestamp := time.Now().Format("20060102_150405")
		*outputFlag = fmt.Sprintf("backup_%s.sql", timestamp)
	}

	fmt.Printf("Creating database backup to %s...\n", *outputFlag)

	// Create backup directory if it doesn't exist
	backupDir := filepath.Dir(*outputFlag)
	if backupDir != "." {
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			log.Fatal("Failed to create backup directory:", err)
		}
	}

	// Use mysqldump command line tool for backup
	cmd := exec.Command("mysqldump",
		"-h", cfg.DBHost,
		"-P", fmt.Sprintf("%d", cfg.DBPort),
		"-u", cfg.DBUser,
		"-p"+cfg.DBPassword,
		"--single-transaction",
		"--routines",
		"--triggers",
		cfg.DBName)

	outputFile, err := os.Create(*outputFlag)
	if err != nil {
		log.Fatal("Failed to create output file:", err)
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to create backup:", err)
	}

	fmt.Printf("Database backup created successfully: %s\n", *outputFlag)
}

// Fix fixes common database issues
func Fix(cfg *config.Config) {
	flagSet := flag.NewFlagSet("fix", flag.ExitOnError)
	checkOnly := flagSet.Bool("check-only", false, "Only check for issues, don't fix them")
	flagSet.Parse(os.Args[2:])

	fmt.Println("Checking database for common issues...")

	db, err := sql.Open("mysql", getDSN(cfg))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 1. 讀取 articles 表
	rows, err := db.Query("SELECT id, path FROM articles")
	if err != nil {
		log.Fatal("Failed to query articles:", err)
	}
	defer rows.Close()

	type Article struct {
		ID   int64
		Path string
	}
	var curArticles []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Path); err != nil {
			log.Printf("Error scanning article: %v\n", err)
			continue
		}
		curArticles = append(curArticles, a)
	}

	// 2. 檢查文件是否存在
	var delArticles []int64
	for _, a := range curArticles {
		articlePath := filepath.Join("articles", a.Path)
		if _, err := os.Stat(articlePath); os.IsNotExist(err) {
			delArticles = append(delArticles, a.ID)
			fmt.Printf("Article file not found: %s (id=%d)\n", articlePath, a.ID)
		}
	}

	// 3. 刪除 DB 中 "現實中 不存在的 article"
	deleteArticleByID(db, delArticles, *checkOnly)

	// 4. Recursive ./articles
	var foundFiles []string
	err = filepath.Walk("articles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel("articles", path)
			if err != nil {
				return err
			}
			foundFiles = append(foundFiles, relPath)
		}
		return nil
	})
	if err != nil && err.Error() != "stop" {
		fmt.Printf("Error walking articles directory: %v\n", err)
	}

	// 用 map 方便查找 DB 中已存在的 path
	curPathMap := make(map[string]struct{})
	for _, a := range curArticles {
		curPathMap[a.Path] = struct{}{}
	}

	// 選擇更新 or 新增
	fmt.Printf("Found %d files in filesystem\n", len(foundFiles))
	fmt.Printf("Found %d paths in database\n", len(curPathMap))

	// Debug: 顯示資料庫中的路徑
	fmt.Println("Database paths:")
	for path := range curPathMap {
		fmt.Printf("  DB: '%s'\n", path)
	}

	// Debug: 顯示文件系統中的路徑
	fmt.Println("Filesystem paths:")
	for _, f := range foundFiles {
		fmt.Printf("  FS: '%s'\n", f)
	}

	for i, f := range foundFiles {
		// print(title, create_date, edit_date)
		if _, exists := curPathMap[f]; exists {
			fmt.Printf("  %d. update %s\n", i+1, f)
		} else {
			fmt.Printf("  %d. insert %s (NOT IN DB)\n", i+1, f)
			// 檢查是否有相似路徑（忽略大小寫）
			for dbPath := range curPathMap {
				if strings.ToLower(f) == strings.ToLower(dbPath) {
					fmt.Printf("     -> Similar path in DB: %s\n", dbPath)
				}
			}
		}
	}

	// Check for orphaned records in TABLE article_tags
	checkOrphanedRecords(db, "article_tags", *checkOnly)

	// Check for orphaned records in TABLE search_index
	checkOrphanedRecords(db, "search_index", *checkOnly)

	// Check for duplicate entries in Table tags
	checkDuplicateEntries(db, *checkOnly)

	if *checkOnly {
		fmt.Println("✅ Database check completed!")
	} else {
		fmt.Println("✅ Database fixes completed!")
	}
}

// Status checks database status
func Status(cfg *config.Config) {
	fmt.Println("Checking database status...")

	db, err := sql.Open("mysql", getDSN(cfg))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	fmt.Println("Database connection successful")

	// Get database info
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Printf("Warning: Failed to get database version: %v\n", err)
	} else {
		fmt.Printf("Database version: %s\n", version)
	}

	// Get table info (name and columns)
	tables := []string{"articles", "tags", "article_tags", "search_index"}
	for _, table := range tables {
		fmt.Printf("\nTable: %s\n", table)

		// Get column information
		rows, err := db.Query(fmt.Sprintf("DESCRIBE %s", table))
		if err != nil {
			fmt.Printf("Error getting columns: %v\n", err)
			continue
		}
		defer rows.Close()

		fmt.Printf("Columns:\n")
		for rows.Next() {
			var field, typ, null, key, defaultVal, extra sql.NullString
			err := rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra)
			if err != nil {
				fmt.Printf("Error reading column info: %v\n", err)
				continue
			}

			// Format the column info
			columnInfo := fmt.Sprintf("    - %s (%s)", field.String, typ.String)
			if key.String == "PRI" {
				columnInfo += " [PRIMARY KEY]"
			}
			if null.String == "NO" {
				columnInfo += " [NOT NULL]"
			}
			if defaultVal.Valid {
				columnInfo += fmt.Sprintf(" [DEFAULT: %s]", defaultVal.String)
			}
			fmt.Println(columnInfo)
		}

		if err = rows.Err(); err != nil {
			fmt.Printf("Error iterating columns: %v\n", err)
		}
	}
}

// Create drops and recreates the entire database
func Restore(cfg *config.Config) {
	flagSet := flag.NewFlagSet("create", flag.ExitOnError)
	forceFlag := flagSet.Bool("force", false, "Force recreation without confirmation")
	migrateFile := flagSet.String("migrate-file", "init", "SQL file name (without .sql extension) to use for migration")
	flagSet.Parse(os.Args[2:])

	fmt.Printf("Database: %s\n", cfg.DBName)
	fmt.Printf("Migration file: db/%s.sql\n", *migrateFile)

	if *forceFlag {
		fmt.Println("WARNING: This will DELETE the entire database and recreate it!")
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" && response != "y" {
			fmt.Println("Cancelled Drop Database.")
		} else {
			// Drop and create database
			dropCreateDB(cfg)
		}
	}

	// Now run migration to create tables and initial data
	fmt.Printf("Running migration using db/%s.sql...\n", *migrateFile)

	// Temporarily modify os.Args to pass the migrate-file parameter to Migrate
	originalArgs := os.Args
	os.Args = append([]string{os.Args[0], "migrate", "--init=" + *migrateFile})
	defer func() { os.Args = originalArgs }()

	Migrate(cfg)

	fmt.Println("Database creation completed successfully!")
}

// Helper functions

// dropCreateDB drops and recreates the database
func dropCreateDB(cfg *config.Config) {
	fmt.Println("Dropping database...")

	// Connect without specifying database
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		log.Fatal("Failed to connect to database server:", err)
	}
	defer db.Close()

	// Drop database if exists
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName))
	if err != nil {
		log.Fatal("Failed to drop database:", err)
	}
	fmt.Printf("Dropped database: %s\n", cfg.DBName)

	// Create new database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName))
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}
	fmt.Printf("Created database: %s\n", cfg.DBName)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getDSN(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
}

func splitSQL(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for _, r := range sql {
		ch := string(r)
		if ch == "'" && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		}
		if ch == "\"" && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}
		current.WriteRune(r)
		if ch == ";" && !inSingleQuote && !inDoubleQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}
	// Add any remaining statement
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}
	return statements
}

func deleteArticleByID(db *sql.DB, IDArray []int64, checkOnly bool) {
	if len(IDArray) > 0 {
		if checkOnly {
			fmt.Printf("Found %d articles in DB but missing files. (No deletion in check-only mode)\n", len(IDArray))
		} else {
			fmt.Printf("Deleting %d articles from DB that have no corresponding file...\n", len(IDArray))
			for _, id := range IDArray {
				_, err := db.Exec("DELETE FROM articles WHERE id = ?", id)
				if err != nil {
					fmt.Printf("Failed to delete article id=%d: %v\n", id, err)
				} else {
					fmt.Printf("Deleted article id=%d\n", id)
				}
			}
		}
	} else {
		fmt.Println("All articles have corresponding files.")
	}
}

func checkOrphanedRecords(db *sql.DB, tableName string, checkOnly bool) {
	var count int
	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s t
		LEFT JOIN articles a ON t.article_id = a.id
		WHERE a.id IS NULL
	`, tableName)

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Printf("⚠️  Error checking orphaned %s: %v\n", tableName, err)
		return
	}

	if count > 0 {
		fmt.Printf("Found %d orphaned %s records\n", count, tableName)
		if !checkOnly {
			delQuery := fmt.Sprintf("DELETE FROM %s WHERE article_id NOT IN (SELECT id FROM articles)", tableName)
			_, err := db.Exec(delQuery)
			if err != nil {
				fmt.Printf("Failed to clean orphaned %s: %v\n", tableName, err)
			} else {
				fmt.Printf("Cleaned %d orphaned %s records\n", count, tableName)
			}
		}
	} else {
		fmt.Printf("No orphaned %s records found\n", tableName)
	}
}

func checkDuplicateEntries(db *sql.DB, checkOnly bool) {
	// Check for duplicate tags
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT name, COUNT(*) as cnt FROM tags GROUP BY name HAVING cnt > 1
		) as duplicates
	`).Scan(&count)

	if err != nil {
		fmt.Printf("⚠️  Error checking duplicate tags: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Printf("⚠️  Found %d duplicate tag names\n", count)
		if !checkOnly {
			// Keep the first occurrence of each tag name
			_, err := db.Exec(`
				DELETE t1 FROM tags t1
				INNER JOIN tags t2 
				WHERE t1.id > t2.id AND t1.name = t2.name
			`)
			if err != nil {
				fmt.Printf("Failed to clean duplicate tags: %v\n", err)
			} else {
				fmt.Printf("Cleaned duplicate tags\n")
			}
		}
	} else {
		fmt.Println("No duplicate tags found")
	}
}
