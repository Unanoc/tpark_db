package helpers

import (
	"io/ioutil"
	"log"
	"sync"
	"tpark_db/database"
	"tpark_db/models"
)

func getCountOfTable(wg *sync.WaitGroup, table string, status *models.Status) {
	defer wg.Done()
	queryString := "SELECT COUNT(*) FROM " + table
	var count int
	rows := database.DB.Conn.QueryRow(queryString)
	_ = rows.Scan(&count)

	switch table {
	case "users":
		status.User = count
	case "forums":
		status.Forum = count
	case "threads":
		status.Thread = count
	case "posts":
		status.Post = count
	default:
		log.Panic("wrong type of table")
	}
}

// StatusHelper selects count of each table.
func StatusHelper() (*models.Status, error) {
	status := models.Status{}
	listTables := []string{"users", "forums", "threads", "posts"}
	wg := &sync.WaitGroup{}
	for _, table := range listTables {
		wg.Add(1)
		go getCountOfTable(wg, table, &status)
	}
	wg.Wait()
	return &status, nil
}

// ExecSQLScript executes sql script.
func ExecSQLScript(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := database.DB.Conn.Exec(string(content)); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
