package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()
	res, err := db.Exec(`INSERT INTO Clients(ip,port,system,client_key,date) VALUES("192.168.0.166","5555","Windows","123456", "14-07-2024");`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}

/*func main() {
	add_client()
}*/
