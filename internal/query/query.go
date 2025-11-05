package query

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type Data struct {
	id       int
	name     string
	username string
}

func GetUserId(userField *Data) int {
	return userField.id
}

func GetUserName(userField *Data) string {
	return userField.name
}

func GetUserUsername(userField *Data) string {
	return userField.username
}

func OneUser(db *sql.DB, id int) (*Data, error) {
	qy := `SELECT * FROM user WHERE id = ?`
	row := db.QueryRow(qy, id)                                  // Crear la query
	result := &Data{}                                           // Estructura para los datos a recibir de la query
	err := row.Scan(&result.id, &result.name, &result.username) // Guardar los datos del resultado de la query en la estructura
	if err != nil {
		fmt.Printf("An error occur \n%s\n", err)
		return nil, err
	}
	return result, nil
}

func ManyQuery(db *sql.DB) ([]Data, error) {
	qy := `SELECT * FROM user ORDER BY id`
	row, err := db.Query(qy)
	if err != nil {
		fmt.Printf("An error occur \n%s\n", err)
		return nil, err
	}
	defer row.Close()
	var result []Data
	for row.Next() {
		temp := &Data{}
		err = row.Scan(&temp.id, &temp.name, &temp.username)
		if err != nil {
			fmt.Printf("An error occur during the scan in the query\n%s\n", err)
			return nil, err
		}
		result = append(result, *temp)
	}
	return result, nil
}
