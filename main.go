package main

import (
	"flag"
	"log"
	"os"
)

var (
	dbaddr     = flag.String("dbhost", os.Getenv("dbaddr"), "IP address for DB.")
	dbport     = flag.Int("dbport", 5432, "Port number for DB.")
	dbuser     = flag.String("dbuser", os.Getenv("dbuser"), "DB username.")
	dbname     = flag.String("dbname", os.Getenv("dbname"), "DB name.")
	dbpassword = flag.String("dbpassword", os.Getenv("dbpassword"), "DB password.")
	addr       = flag.String("addr", "localhost:8101", "The address to connect to")
)

func main() {
	flag.Parse()

	err := client()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
