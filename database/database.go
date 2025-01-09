package database

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %q", err)
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}

func EnsureTablesExist(db *sql.DB) error {
	// Define the names of the tables to be created
	tableNames := []string{"users", "offers", "orders"}

	for _, tableName := range tableNames {

		// Check if the table exists
		query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
		var exists bool
		err := db.QueryRow(query, tableName).Scan(&exists)
		if err != nil {
			return err
		}

		// If the table exists, log a message
		if exists {
			log.Printf("Tabla '%s' ya existe - OK!", tableName)
		} else {

			// if not exists, create the table
			createTableQuery := getCreateTableQuery(tableName)
			_, err := db.Exec(createTableQuery)
			if err != nil {
				return err
			}
			if tableName == "offers" {
				err = insertDefaultOffers(db)
				if err != nil {
					return err
				}
			}
			log.Printf("Tabla '%s' creada correctamente", tableName)
		}
	}

	return nil
}

// Obtains the query to create the table
func getCreateTableQuery(tableName string) string {
	switch tableName {
	case "users":
		return `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL
        ); 
        INSERT INTO users (username, email, password)
        VALUES 
        ('Ubuntu', 'Ubuntu@gmail.com', 'Ubuntu')
        ON CONFLICT (email) DO NOTHING;`

	case "offers":
		return `CREATE TABLE IF NOT EXISTS offers (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL,
            price FLOAT NOT NULL,
            category VARCHAR(50) NOT NULL
        );`

	case "orders":
		return `CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            status VARCHAR(50) NOT NULL,
            user_id INT NOT NULL,
			total INT NOT NULL,
            CONSTRAINT fk_user
                FOREIGN KEY(user_id) 
	            REFERENCES users(id)
        );`

	default:
		return "" // Return empty string if the table name is not found
	}
}

// function to preload the prices of the offers
func insertDefaultOffers(db *sql.DB) error {
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)

	offers := []struct {
		Name     string
		Category string
	}{
		{"meat", "food"},
		{"vegetables", "food"},
		{"fruits", "food"},
		{"water", "drink"},
		{"antibiotics", "medicine"},
		{"analgesics", "medicine"},
		{"bandages", "medicine"},
		{"pistol ammo", "ammo"},
		{"rifle ammo", "ammo"},
		{"shotgun ammo", "ammo"},
	}

	for _, offer := range offers {
		price := float64(random.Intn(20) + 1) // Generates a random price between 1 and 20
		query := `INSERT INTO offers (name, price, category) VALUES ($1, $2, $3)`
		_, err := db.Exec(query, offer.Name, price, offer.Category)
		if err != nil {
			return err
		}
	}

	log.Println("Default offers inserted successfully")
	return nil
}
