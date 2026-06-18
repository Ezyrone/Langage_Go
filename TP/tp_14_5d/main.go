package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `db:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func CreateUserSQL(db *sql.DB, user User) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetUsersSQL(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetUserByIDSQL(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserSQL(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func DeleteUserSQL(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CreateUserSQLX(db *sqlx.DB, user User) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetUsersSQLX(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT id, name, email FROM users")
	return users, err
}

func GetUserByIDSQLX(db *sqlx.DB, id int) (*User, error) {
	var u User
	err := db.Get(&u, "SELECT id, name, email FROM users WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserSQLX(db *sqlx.DB, user User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func DeleteUserSQLX(db *sqlx.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CreateUserGORM(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUsersGORM(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserByIDGORM(db *gorm.DB, id int) (*User, error) {
	var u User
	err := db.First(&u, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserGORM(db *gorm.DB, user User) error {
	return db.Save(&user).Error
}

func DeleteUserGORM(db *gorm.DB, id int) error {
	return db.Delete(&User{}, id).Error
}

func main() {
	fmt.Println("=== Exercice 1.1 & 1.2 : database/sql ===")

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connecté à la base de données SQLite (database/sql).")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	id1, err := CreateUserSQL(db, User{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[SQL] Utilisateur créé avec ID: %d\n", id1)

	id2, err := CreateUserSQL(db, User{Name: "Bob", Email: "bob@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[SQL] Utilisateur créé avec ID: %d\n", id2)

	users, err := GetUsersSQL(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQL] Tous les utilisateurs:")
	for _, u := range users {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	user, err := GetUserByIDSQL(db, int(id1))
	if err != nil {
		log.Fatal(err)
	}
	if user != nil {
		fmt.Printf("[SQL] Utilisateur ID=%d: Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
	}

	err = UpdateUserSQL(db, User{ID: int(id1), Name: "Alice Updated", Email: "alice.updated@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQL] Utilisateur mis à jour.")

	user, err = GetUserByIDSQL(db, int(id1))
	if err != nil {
		log.Fatal(err)
	}
	if user != nil {
		fmt.Printf("[SQL] Après mise à jour: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
	}

	err = DeleteUserSQL(db, int(id2))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQL] Utilisateur supprimé.")

	users, err = GetUsersSQL(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQL] Utilisateurs restants:")
	for _, u := range users {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	fmt.Println("\n=== Exercice 1.3 : sqlx ===")

	dbx, err := sqlx.Open("sqlite3", "./test_sqlx.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	if err := dbx.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connecté à la base de données SQLite (sqlx).")

	_, err = dbx.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	id1, err = CreateUserSQLX(dbx, User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[SQLX] Utilisateur créé avec ID: %d\n", id1)

	id2, err = CreateUserSQLX(dbx, User{Name: "Diana", Email: "diana@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[SQLX] Utilisateur créé avec ID: %d\n", id2)

	usersx, err := GetUsersSQLX(dbx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQLX] Tous les utilisateurs:")
	for _, u := range usersx {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	userx, err := GetUserByIDSQLX(dbx, int(id1))
	if err != nil {
		log.Fatal(err)
	}
	if userx != nil {
		fmt.Printf("[SQLX] Utilisateur ID=%d: Name=%s, Email=%s\n", userx.ID, userx.Name, userx.Email)
	}

	err = UpdateUserSQLX(dbx, User{ID: int(id1), Name: "Charlie Updated", Email: "charlie.updated@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQLX] Utilisateur mis à jour.")

	userx, err = GetUserByIDSQLX(dbx, int(id1))
	if err != nil {
		log.Fatal(err)
	}
	if userx != nil {
		fmt.Printf("[SQLX] Après mise à jour: ID=%d, Name=%s, Email=%s\n", userx.ID, userx.Name, userx.Email)
	}

	err = DeleteUserSQLX(dbx, int(id2))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQLX] Utilisateur supprimé.")

	usersx, err = GetUsersSQLX(dbx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[SQLX] Utilisateurs restants:")
	for _, u := range usersx {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	fmt.Println("\n=== Exercice 2 : GORM ===")

	gormDB, err := gorm.Open(sqlite.Open("gorm_test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Échec de la connexion GORM:", err)
	}
	log.Println("Connecté à la base de données SQLite (GORM).")

	if err := gormDB.AutoMigrate(&User{}); err != nil {
		log.Fatal("Échec de l'auto-migration GORM:", err)
	}
	log.Println("Auto-migration GORM effectuée.")

	gormUser1 := User{Name: "Eve", Email: "eve@example.com"}
	if err := CreateUserGORM(gormDB, &gormUser1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[GORM] Utilisateur créé avec ID: %d\n", gormUser1.ID)

	gormUser2 := User{Name: "Frank", Email: "frank@example.com"}
	if err := CreateUserGORM(gormDB, &gormUser2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[GORM] Utilisateur créé avec ID: %d\n", gormUser2.ID)

	gormUsers, err := GetUsersGORM(gormDB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[GORM] Tous les utilisateurs:")
	for _, u := range gormUsers {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}

	gormUser, err := GetUserByIDGORM(gormDB, gormUser1.ID)
	if err != nil {
		log.Fatal(err)
	}
	if gormUser != nil {
		fmt.Printf("[GORM] Utilisateur ID=%d: Name=%s, Email=%s\n", gormUser.ID, gormUser.Name, gormUser.Email)
	}

	err = UpdateUserGORM(gormDB, User{ID: gormUser1.ID, Name: "Eve Updated", Email: "eve.updated@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[GORM] Utilisateur mis à jour.")

	gormUser, err = GetUserByIDGORM(gormDB, gormUser1.ID)
	if err != nil {
		log.Fatal(err)
	}
	if gormUser != nil {
		fmt.Printf("[GORM] Après mise à jour: ID=%d, Name=%s, Email=%s\n", gormUser.ID, gormUser.Name, gormUser.Email)
	}

	err = DeleteUserGORM(gormDB, gormUser2.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[GORM] Utilisateur supprimé.")

	gormUsers, err = GetUsersGORM(gormDB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[GORM] Utilisateurs restants:")
	for _, u := range gormUsers {
		fmt.Printf("  ID=%d, Name=%s, Email=%s\n", u.ID, u.Name, u.Email)
	}
}
