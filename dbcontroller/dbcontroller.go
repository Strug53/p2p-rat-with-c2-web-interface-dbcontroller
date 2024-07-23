package dbcontroller

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	Id         int
	IP         string
	Port       string
	System     string
	Client_key string
	Date       string
}

func Delete_from_id(id int) {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()

	res, err := db.Exec(
		"DELETE FROM Clients WHERE id = $1;UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'Clients';", id)
	if err != nil {
		fmt.Println(err)
	}
	_ = res
}
func Delete_from_IP(ip string) {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()

	res, err := db.Exec(
		"DELETE FROM Clients WHERE ip = $1;UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'Clients';", ip)
	if err != nil {
		fmt.Println(err)
	}
	_ = res

}

func Delete_all() {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()

	res, err := db.Exec(
		"DELETE FROM Clients; UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'Clients';")
	if err != nil {
		fmt.Println(err)
	}
	_ = res
}

func Add_client(ip string, port string, system string, client_key string, date string) {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()
	res, err := db.Exec(`INSERT INTO Clients(ip,port,system,client_key,date) VALUES($1,$2,$3,$4,$5);`, ip, port, system, client_key, date)
	if err != nil {
		fmt.Println(err)
	}
	_ = res

}
func Select_all_clients() []Client {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Clients")
	if err != nil {
		fmt.Println(err)
	}
	clients := []Client{}
	defer rows.Close()

	for rows.Next() {
		c := Client{}
		err := rows.Scan(&c.Id, &c.IP, &c.Port, &c.System, &c.Client_key, &c.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clients = append(clients, c)
	}
	return clients
}
func Select_client_ID(uid int) Client {
	db, err := sql.Open("sqlite3", "clients.db")
	if err != nil {
		fmt.Printf("Error in opening db \n")
	}
	fmt.Printf("Opened")
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM Clients WHERE id = %s", strconv.Itoa(uid)))
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	c := Client{}
	for rows.Next() {
		err := rows.Scan(&c.Id, &c.IP, &c.Port, &c.System, &c.Client_key, &c.Date)
		if err != nil {
			fmt.Println(err)
			return Client{}
		}
	}

	return c
}

// UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'ranobes';
// func main() {

// }
