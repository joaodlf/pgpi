package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	defaultPgHost = "localhost"
	defaultPgPort = "5432"
	defaultPgUser = "postgres"
	green         *color.Color
	yellow        *color.Color
	red           *color.Color
	cyan          *color.Color
)

func init() {
	// Exit execution for certain syscalls.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _ = range c {
			os.Exit(0)
		}
	}()

	// Load the different colors to be used in the CLI.
	green = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
	red = color.New(color.FgRed)
	cyan = color.New(color.FgCyan)
}

// validateConfirmation validates confirmation ([y/n]) user input.
func validateConfirmation(value string) (bool, error) {
	if value == "y" || value == "Y" {
		return true, nil
	}

	if value == "n" || value == "N" {
		return false, nil
	}

	return false, errors.New("invalid value supplied")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// PG HOST.
	green.Print("Postgres Host (defaults to 'localhost'): ")
	scanner.Scan()
	pgHost := scanner.Text()

	if pgHost == "" {
		pgHost = defaultPgHost
	}

	// PG PORT.
	green.Print("Postgres Port (defaults to 5432): ")
	scanner.Scan()
	pgPort := scanner.Text()

	if pgPort == "" {
		pgPort = defaultPgPort
	}
	pgPortInt, _ := strconv.Atoi(pgPort)

	// PG USER.
	green.Print("Postgres User (defaults to 'postgres'): ")
	scanner.Scan()
	pgUser := scanner.Text()

	if pgUser == "" {
		pgUser = defaultPgUser
	}

	// PG USER PASSWORD.
	passwordPrompt := fmt.Sprintf("Password for user '%s': ", pgUser)
	pgPassword := ""
	for {
		green.Print(passwordPrompt)
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		pgPassword = string(bytePassword)
		fmt.Println("")

		if pgPassword != "" {
			break
		}
	}

	// PG DATABASE.
	pgDatabase := ""
	for {
		green.Print("Postgres Database: ")
		scanner.Scan()
		pgDatabase = scanner.Text()

		if pgDatabase != "" {
			break
		}
	}

	// DB Connection.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPortInt, pgUser, pgPassword, pgDatabase)

	db, err := sql.Open("postgres", psqlInfo)

	yellow.Println("Connecting to Postgres database...")

	if err != nil {
		red.Println(err)
		os.Exit(0)
	}

	err = db.Ping()

	if err != nil {
		red.Println(err)
		os.Exit(0)
	}

	yellow.Println("Connnected!")

tr:
	// TABLE NAME REGEX.
	tableNameRegex := ""
	for {
		green.Print("Table name regex to apply index (E.g. tablename_.*_.*): ")
		scanner.Scan()
		tableNameRegex = scanner.Text()

		if tableNameRegex != "" {
			break
		}
	}

	tablesQuery := fmt.Sprintf(`SELECT tablename 
								FROM pg_tables 
								WHERE SUBSTRING(tablename FROM '%s') <> '';`, tableNameRegex)

	yellow.Println("Finding tables/partitions...")
	rows, _ := db.Query(tablesQuery)

	var tables []string
	for rows.Next() {
		var tablename string
		if err := rows.Scan(&tablename); err != nil {
			fmt.Println(err)
		}
		cyan.Printf("%s\n", tablename)
		tables = append(tables, tablename)
	}

	var confirmation bool
	for {
		green.Print("Is this correct? [y/n] ")
		scanner.Scan()
		confirmation, err = validateConfirmation(scanner.Text())

		if err == nil {
			break
		}
	}

	if confirmation == false {
		goto tr
	}

	// INDEX NAME.
	indexName := ""
	for {
		green.Print("New index name: ")
		scanner.Scan()
		indexName = scanner.Text()

		if indexName != "" {
			break
		}
	}

	// UNIQUE INDEX.
	var uniqueIndex bool
	for {
		green.Print("Is this a unique index? [y/n]: ")
		scanner.Scan()
		uniqueIndex, err = validateConfirmation(scanner.Text())

		if err == nil {
			break
		}
	}

	var uniqueIndexStr string
	if uniqueIndex {
		uniqueIndexStr = "UNIQUE"
	}

	indexColumns := ""
	for {
		green.Print("Index columns (E.g. col1, col2 DESC, col3): ")
		scanner.Scan()
		indexColumns = scanner.Text()

		if indexColumns != "" {
			break
		}
	}

	var queries []string
	for _, table := range tables {
		fullIndexName := fmt.Sprintf("%s_%s", table, indexName)
		indexQuery := fmt.Sprintf(`CREATE %s INDEX CONCURRENTLY %s ON %s (%s);`, uniqueIndexStr, fullIndexName, table, indexColumns)
		queries = append(queries, indexQuery)
		cyan.Println(indexQuery)
	}

	var execute bool
	for {
		green.Print("FINAL STEP! Execute above queries? [y/n]: ")
		scanner.Scan()
		execute, err = validateConfirmation(scanner.Text())

		if err == nil {
			break
		}
	}

	if !execute {
		os.Exit(0)
	}
	for _, query := range queries {
		cyan.Println(fmt.Sprintf("Executing '%s'", query))
		_, err := db.Exec(query)

		if err != nil {
			red.Println(err)
		}
	}

	green.Println("All queries executed!")
}
