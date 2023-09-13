package models

import (
	"database/sql"
	"fmt"
	"strconv"
)

func UpdateTable(db *sql.DB, table string, id int, setAttrs map[string]string) error {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range setAttrs {
		queryStr += fmt.Sprintf("%s=$%d,", key, i)
		values = append(values, val)
		i++
	}
	queryStr = queryStr[:len(queryStr)-1]
	values = append(values, id)
	_, err := db.Exec("UPDATE "+table+" SET "+queryStr+" WHERE id=$"+strconv.Itoa(i), values...)
	if err != nil {
		return err
	}
	return nil
}
