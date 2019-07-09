package main

import "fmt"
import "github.com/gocql/gocql"

var Session *gocql.Session

type User struct {
	id        string
	firstName string
	lastName  string
	age       int
}

func Init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "exampleks"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra Session created")
}

func creatUser(usr User) {
	fmt.Println("Creating new user\n", usr)
	err := Session.Query("INSERT INTO users(usrid, first_name, last_name, age) VALUES(?, ?, ?, ?)",
		usr.id, usr.firstName, usr.lastName, usr.age).Exec()
	if err != nil {
		fmt.Println("Error while inserting new user")
		fmt.Println(err)
	}
}

func getUsers() []User {
	fmt.Println("Getting all Users")
	var users []User
	m := map[string]interface{}{}

	iter := Session.Query("SELECT * FROM users").Iter()
	for iter.MapScan(m) {
		users = append(users, User{
			id:        m["usrid"].(string),
			firstName: m["first_name"].(string),
			lastName:  m["last_name"].(string),
			age:       m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	return users
}

func updateUser(usr User) {
	fmt.Printf("Updating User with id = %s\n", usr.id)
	if err := Session.Query("UPDATE users SET first_name = ?, last_name = ?, age = ? WHERE usrid = ?",
		usr.firstName, usr.lastName, usr.age, usr.id).Exec(); err != nil {
		fmt.Println("Error while updating User")
		fmt.Println(err)
	}
}

func deleteUser(id string) {
	fmt.Printf("Deleting User with id = %s\n", id)
	if err := Session.Query("DELETE FROM users WHERE usrid = ?", id).Exec(); err != nil {
		fmt.Println("Error while deleting User")
		fmt.Println(err)
	}
}

func main() {
	Init()
	creatUser(User{"User-Id_1", "Mac", "Davis", 20})
	creatUser(User{"User-Id_2", "John", "Truman", 20})
	creatUser(User{"User-Id_3", "Rob", "Freeman", 20})
	updateUser(User{"User-Id_1", "Micheal", "Davis", 20})
	fmt.Println(getUsers())
	deleteUser("User-Id_2")
	fmt.Println(getUsers())
}
