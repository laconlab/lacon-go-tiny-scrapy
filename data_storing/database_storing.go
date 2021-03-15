package data_storing

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"lacon-go-tiny-scrapy/logger"
	"lacon-go-tiny-scrapy/parser"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	host string
	username string
	password string
	databaseName string
	numberOfWorkers int

	workersGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "database_worker_count",
		Help: "Number of active database workers",
	})

	databaseSummary = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "database_worker_summary",
		Help: "Performances per worker",
	})

	errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "database_error_count",
		Help: "Number of errors in database",
	})
)

func init() {
	host = os.Getenv("DB_HOST")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	databaseName = os.Getenv("DB_NAME")
	logger.INFO.Println("Creating table")
	createTable()
	logger.INFO.Println("Table created")

	config := newDatabaseConfiguration()
	numberOfWorkers = config.getWorkerPoolSize()
}

func createTable() {

	query := getCreateTableQuery()
	db := openConnection()
	defer db.Close()

	_, err := db.Exec(query)
	if err != nil {
		logger.ERROR.Printf("Error while creating table: %v", err)
		errorCounter.Inc()
		panic(err)
	}
}

func getCreateTableQuery() string {
	query := "CREATE TABLE IF NOT EXISTS %s.raw_text (" +
		"id INT, website VARCHAR(50), publish_date varchar(100)," +
		"dt DATETIME DEFAULT CURRENT_TIMESTAMP, url varchar(1000), title varchar(1000)," +
		"article MEDIUMTEXT)" +
		"ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8mb4 " +
		"COLLATE=utf8mb4_unicode_ci;"

	return fmt.Sprintf(query, databaseName)
}

func NewDatabaseInsertWorker(wg *sync.WaitGroup, inputData <-chan parser.ParsedContent) {
	wg.Add(numberOfWorkers)
	logger.INFO.Printf("Starting %d database workers\n", numberOfWorkers)
	workersGauge.Set(float64(numberOfWorkers))
	for i := 0; i < numberOfWorkers; i++ {
		go startInsertingInDatabase(wg, inputData)
	}
}

func startInsertingInDatabase(wg *sync.WaitGroup, inputData <-chan parser.ParsedContent) {
	db := openConnection()
	defer db.Close()
	for content := range inputData {
		startTime := time.Now().UnixNano()
		insert(db,content)
		endTime := time.Now().UnixNano()
		databaseSummary.Observe(float64(endTime - startTime))
	}
	wg.Done()
}

func insert(db *sql.DB, content parser.ParsedContent) {
	names, values := content.GetFieldNamesAndValues()
	query := createQuery(names)
	tryInsert(db, query, values)
}

func tryInsert(db *sql.DB, query string, values []interface{}) {
	q, err := db.Query(query, values...)
	if err != nil {
		logger.ERROR.Printf("Error inserting into database: %v\n", err)
		logger.ERROR.Printf("Failed to insert %v\n", values)
		errorCounter.Inc()
		return
	}
	_ = q.Close()
}

func createQuery(fieldNames []string) string {
	queryTemplate := "INSERT INTO %s.raw_text (%s) VALUES (%s)"
	names := strings.Join(fieldNames, ",")
	questionMarks := createNQuestionMarks(len(fieldNames))
	return fmt.Sprintf(queryTemplate, databaseName, names, questionMarks)
}

func openConnection() *sql.DB {
	dataSourceName := createDataSourceName()
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.ERROR.Printf("Error whit opening database connection %v\n", err)
		errorCounter.Inc()
	}
	return db
}

func createDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, databaseName)
}

func createNQuestionMarks(n int) string {
	q := strings.Repeat("?,", n)
	return q[:len(q)-1]
}