package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type service struct {
	db *pgx.Conn
}

func main() {
	godotenv.Overload()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	serv := service{
		db: conn,
	}

	serv.runMenu()
}

func (s *service) runMenu() {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========== SQL Injection Demo ==========")
		fmt.Println("1. Vulnerable login (or '1'='1')")
		fmt.Println("2. Secure login (parameterized queries)")
		fmt.Println("3. Vulnerable user search (union injection)")
		fmt.Println("4. Secure user search (parameterized queries)")
		fmt.Println("5. Vulnerable user delete ()")
		fmt.Println("6. Secure user delete (parameterized queries)")
		fmt.Println("7. Exit")
		fmt.Print("\nChoose option: ")

		choice, _ := r.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			s.vulnerableLogin(r)
		case "2":
			s.secureLogin(r)
		case "3":
			s.vulnerableSearchUser(r)
		case "4":
			s.secureUserSearch(r)
		case "5":
			s.vulnarableDeleteUser(r)
		case "6":
			s.secureDeleteUser(r)
		case "7":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option, try again")
		}
	}
}

func (s *service) vulnerableSearchUser(r *bufio.Reader) {
	fmt.Println("\n--- Vulnerable Search (NO PROTECTION) ---")
	fmt.Print("Enter search term: ")

	searchTerm, _ := r.ReadString('\n')
	searchTerm = strings.TrimSpace(searchTerm)

	query := fmt.Sprintf("select username, role from users where username like '%%%s%%'", searchTerm)
	fmt.Printf("[DEBUG] Query: %s\n", query)

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("[FAIL] Search failed: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var role string
		err := rows.Scan(&username, &role)
		if err != nil {
			fmt.Printf("[FAIL] Unable to scan row: %v\n", err)
			return
		}
		fmt.Printf("%s, %s\n", username, role)
	}
}

func (s *service) secureUserSearch(r *bufio.Reader) {
	fmt.Println("\n--- Secure user search (PROTECTED) ---")
	fmt.Print("Enter search term: ")
	searchTerm, _ := r.ReadString('\n')
	searchTerm = "%" + strings.TrimSpace(searchTerm) + "%"

	query := "select username, role from users where username like $1"
	fmt.Printf("[DEBUG] Query: %s\n", query)
	fmt.Printf("[DEBUG] Parameters: search term=%q\n", searchTerm)

	rows, err := s.db.Query(context.Background(), query, searchTerm)
	if err != nil {
		fmt.Printf("[FAIL] Search failed: %v\n", err)
		return
	}

	for rows.Next() {
		var username string
		var role string
		err := rows.Scan(&username, &role)
		if err != nil {
			fmt.Printf("[FAIL] Unable to scan row: %v\n", err)
			return
		}
		fmt.Printf("%s, %s\n", username, role)
	}
}

func (s *service) vulnerableLogin(r *bufio.Reader) {
	fmt.Println("\n--- Vulnerable Login (NO PROTECTION) ---")
	fmt.Print("Enter username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := r.ReadString('\n')
	password = strings.TrimSpace(password)

	query := fmt.Sprintf("SELECT id, username, password, role FROM users WHERE username='%s' AND password='%s'", username, password)
	fmt.Printf("[DEBUG] Query: %s\n", query)

	var id int
	var dbUsername, dbPassword, role string

	err := s.db.QueryRow(context.Background(), query).Scan(&id, &dbUsername, &dbPassword, &role)
	if err != nil {
		fmt.Printf("[FAIL] Authentication failed: %v\n", err)
		return
	}

	fmt.Printf("[SUCCESS] User authenticated!\n")
	fmt.Printf("  ID: %d\n", id)
	fmt.Printf("  Username: %s\n", dbUsername)
	fmt.Printf("  Password: %s (has been stolen)\n", dbPassword)
	fmt.Printf("  Role: %s\n", role)
}

func (s *service) secureLogin(r *bufio.Reader) {
	fmt.Println("\n--- Secure Login (PROTECTED) ---")
	fmt.Print("Enter username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := r.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Println("\n[SECURE] Executing parameterized query...")

	query := "SELECT id, username, password, role FROM users WHERE username=$1 AND password=$2"
	fmt.Printf("[DEBUG] Query: %s\n", query)
	fmt.Printf("[DEBUG] Parameters: username=%q, password=%q\n", username, password)

	var id int
	var dbUsername, dbPassword, role string

	err := s.db.QueryRow(context.Background(), query, username, password).Scan(&id, &dbUsername, &dbPassword, &role)
	if err != nil {
		fmt.Printf("[FAIL] Authentication failed: %v\n", err)
		return
	}

	fmt.Printf("[SUCCESS] User authenticated!\n")
	fmt.Printf("  ID: %d\n", id)
	fmt.Printf("  Username: %s\n", dbUsername)
	fmt.Printf("  Role: %s\n", role)
}

func (s *service) vulnarableDeleteUser(r *bufio.Reader) {
	fmt.Println("\n--- Vulnerable User Delete (NO PROTECTION) ---")
	fmt.Print("Enter username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	query := fmt.Sprintf("delete from users where username='%s'", username)
	fmt.Printf("[DEBUG] Query: %s\n", query)

	commandTag, err := s.db.Exec(context.Background(), query)
	if err != nil {
		fmt.Printf("[FAIL] Deletion failed: %v\n", err)
		return
	}

	if commandTag.RowsAffected() == 0 {
		fmt.Println("[WARNING] No rows were affected")
	} else {
		fmt.Printf("User successfully deleted")
	}
}

func (s *service) secureDeleteUser(r *bufio.Reader) {
	fmt.Println("\n--- Secure User Delete (PROTECTED) ---")
	fmt.Print("Enter username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)

	query := "delete from users where username=$1"

	fmt.Printf("[DEBUG] Query: %s\n", query)
	fmt.Printf("[DEBUG] Parameters: username=%s\n", username)

	commandTag, err := s.db.Exec(context.Background(), query, username)
	if err != nil {
		fmt.Printf("[FAIL] Deletion failed: %v\n", err)
		return
	}

	if commandTag.RowsAffected() == 0 {
		fmt.Printf("[WARNING] No rows were affected: %v\n", err)
	} else {
		fmt.Printf("User successfully deleted")
	}
}
