# PKMS Database CLI Tool

A command-line interface for managing the PKMS database operations including migration, backup, restore, and maintenance.

## Prerequisites

- Go 1.16 or higher
- MySQL client tools (`mysql` and `mysqldump` commands)
- Access to the MySQL database

## Installation

The CLI tool is part of the backend project. Make sure you have the required dependencies:

```bash
cd backend
go mod tidy
```

## Usage

### Basic Commands

```bash
docker compose up mysql -d
/* 等到 MySQL 完全啟動 */
docker compose ps

cd backend
go run cli/main.go status
go run cli/main.go <command> [options]
```

## Available Commands
### 1. [Migrate](#1-migrate): 執行 sql file
### 2. [Backup](#2-backup): 將 database 資料 dump 到 `backup_timestamp.sql`
### 3. [Restore](#3-restore)
### 4. [Fix](#4-fix)
### 5. [Status](#5-status)
### 6. [Help](#6-help)

------

### 1. Migrate
執行 sql file<br>
Runs database migrations using SQL files from the `db/` directory.

```bash
# Use default init.sql file
go run cli/main.go migrate

# Use specific SQL file
go run cli/main.go migrate --init=init
go run cli/main.go migrate --init=update
```

This will:
- Connect to the database using environment variables
- Execute all SQL statements from the specified file (e.g., `db/init.sql`, `db/update.sql`)
- Create tables, indexes, and initial data if they don't exist
- Default file is `db/init.sql` if no `--init` flag is provided

### 2. Backup
執行 mysqldump 命令，將 DB 導出到 SQL 文件。
1. 會 `CREATE TABLE`
2. 會 `INSERT DATA`
3. 不會 `DROP TABLE`

```bash
# Create backup with default timestamp filename
go run cli/main.go backup

# Create backup with custom filename
go run cli/main.go backup --output=my_backup.sql

# Create backup in specific directory
go run cli/main.go backup --output=backups/backup_$(date +%Y%m%d).sql
```

### 3. Restore
1. Restores database from a backup file.
2. `--force`: Drops and recreates the entire database from scratch.

|flag||
|----|---|
|--force|執行 `DROP Database`|
|--migrate-file='sqlName'|使用 `${sqlName}.sql`|

```bash
# Create database with confirmation prompt
go run cli/main.go restore

# Force create without confirmation
go run cli/main.go restore --force
```

**⚠️ Warning**: This may completely delete the existing database and all its data!

This will:
- Drop the existing database if it exists
- Create a new database with proper charset and collation
- Run migrations to create tables and initial data

### 4. Fix
Fixes common database issues.

```bash
# Check for issues without fixing them
go run cli/main.go fix --check-only

# Fix issues automatically
go run cli/main.go fix
```

<p style="background-color:#ef444433;color:#F87171;">
還沒做 Fix 中的 update & insert
</p>

This will:
- Remove orphaned records in junction tables
- Clean up duplicate tag entries
- Fix invalid foreign key references

### 5. Status
Checks database status and health.

```bash
go run cli/main.go status
```

This will:
- Test database connectivity
- Show database version
- Display record counts for all tables
- Check for common issues

### 6. Help
Shows usage information.

```bash
go run cli/main.go help
```

## Environment Variables

The CLI uses the same environment variables as the main application:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 3306)
- `DB_USER` - Database username (default: root)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: pkms)

## Examples

### Complete Database Setup

```bash
# 1. Check if database is accessible
go run cli/main.go status

# 2. Run migrations to create tables and initial data
go run cli/main.go migrate

# 3. Verify everything is working
go run cli/main.go status
```

### Fresh Database Creation

```bash
# 1. Drop and recreate the entire database
go run cli/main.go create

# 2. Verify the new database is working
go run cli/main.go status
```

### Database Maintenance

```bash
# 1. Create a backup before maintenance
go run cli/main.go backup --output=backup_before_maintenance.sql

# 2. Check for issues
go run cli/main.go fix --check-only

# 3. Fix issues if found
go run cli/main.go fix

# 4. Verify fixes
go run cli/main.go status
```

### Disaster Recovery

```bash
# 1. Restore from backup
go run cli/main.go restore --file=backup_20241201_143022.sql

# 2. Verify restoration
go run cli/main.go status

# 3. Fix any issues that might have occurred
go run cli/main.go fix
```

## Troubleshooting

### Common Issues

1. **"command not found: mysql"**
   - Install MySQL client tools
   - On macOS: `brew install mysql-client`
   - On Ubuntu: `sudo apt-get install mysql-client`

2. **"Access denied" errors**
   - Check your database credentials in environment variables
   - Ensure the user has proper permissions

3. **"Failed to read init.sql"**
   - Make sure you're running the command from the backend directory
   - Verify that `db/init.sql` exists

4. **"Failed to create backup"**
   - Check if the output directory exists and is writable
   - Ensure you have sufficient disk space

### Debug Mode

For more verbose output, you can set the `DEBUG` environment variable:

```bash
DEBUG=1 go run cli/main.go status
```

## File Structure

```
backend/
├── cli/
│   ├── main.go          # CLI entry point
│   ├── commands/
│   │   └── commands.go  # Command implementations
│   └── README.md        # This file
├── db/
│   └── init.sql         # Database schema and initial data
└── ...
```

## Contributing

When adding new commands:

1. Add the command to the switch statement in `cli/main.go`
2. Implement the command function in `cli/commands/commands.go`
3. Update this README with usage examples
4. Test thoroughly with different scenarios 