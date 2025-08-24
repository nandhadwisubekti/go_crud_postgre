package database

import (
	"database/sql"
	"fmt"
	"log"

	"go-crud-employee/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established successfully")

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// CreateTables creates the necessary tables for the application
func (db *DB) CreateTables() error {
	// Create users table for authentication
	userTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create employees table
	employeeTableQuery := `
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		nip VARCHAR(20) UNIQUE NOT NULL,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(20),
		position VARCHAR(100) NOT NULL,
		department VARCHAR(100) NOT NULL,
		salary DECIMAL(15,2),
		hire_date DATE NOT NULL,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes for better performance
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS idx_employees_nip ON employees(nip);",
		"CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);",
		"CREATE INDEX IF NOT EXISTS idx_employees_department ON employees(department);",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
	}

	// Execute table creation queries
	queries := []string{userTableQuery, employeeTableQuery}
	queries = append(queries, indexQueries...)

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// CreateDefaultUser creates a default admin user for testing
func (db *DB) CreateDefaultUser() error {
	// Check if default user already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", "admin").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %v", err)
	}

	if count > 0 {
		log.Println("Default admin user already exists")
		return nil
	}

	// Create default user with hashed password
	// Password: "admin123" (you should change this in production)
	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi" // bcrypt hash of "admin123"

	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3)`

	_, err = db.Exec(query, "admin", "admin@company.com", hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to create default user: %v", err)
	}

	log.Println("Default admin user created successfully (username: admin, password: admin123)")
	return nil
}