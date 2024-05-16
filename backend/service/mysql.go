package service

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

const itemsPerPage = 30

type MonitorJob struct {
	ID      int64  `json:"id"`
	Project string `json:"project"`
	Cron    string `json:"cron"`
	Enable  bool   `json:"enable"`
}



type MonitorMetric struct {
	ID            int64     `json:"id"`
	Project       string    `json:"project"`
	Catalog       string    `json:"catalog"`
	ItemDesc      string    `json:"item_desc"`
	ItemCondition string    `json:"item_condition"`
	ConnectionName string   `json:"conn_name"`
	DashboardURL  string    `json:"dashboard_url"`
	Status        bool      `json:"status"`
	StatusDesc    string    `json:"status_desc"`
	Screen        string    `json:"screen"`
	CheckDate     time.Time `json:"check_date"`
}

type MonitorConnection struct {
    ID        int64  `json:"id"`
    Name      string `json:"conn_name"`
    Username  string `json:"conn_username"`
    Password  string `json:"conn_password"`
    URL       string `json:"conn_url"`
}

func InitDBConnection() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	err = ConnectDB(cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	if err != nil {
		panic(err)
	}

}

// ConnectDB initializes the database connection
func ConnectDB(databaseName, user, password, host, port string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databaseName)
	var err error
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		return err
	}

	_, err = DB.Exec("select 1")
	if err != nil {
		return err
	}
	return nil
}

// InsertMonitorMetric inserts a new monitor metric into the database
func InsertMonitorMetric(item MonitorMetric) (int64, error) {

	result, err := DB.Exec("INSERT INTO monitor_metrics (project, catalog, item_desc, item_condition,conn_name, dashboard_url, status, status_desc, screen) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.Project, item.Catalog, item.ItemDesc, item.ItemCondition,item.ConnectionName, item.DashboardURL, item.Status, item.StatusDesc, item.Screen)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteMonitorMetric deletes a monitor item from the database
func DeleteMonitorMetric(id int64) error {
	_, err := DB.Exec("DELETE FROM monitor_metrics WHERE id = ?", id)
	return err
}

// UpdateMonitorItem updates a monitor item in the database
func UpdateMonitorMetric(item MonitorMetric) error {
	item.CheckDate = time.Now()

	_, err := DB.Exec("UPDATE monitor_metrics SET project = ?, catalog = ?, item_desc = ?, item_condition = ?,conn_name=?, dashboard_url = ?, status = ?, status_desc = ?, screen = ?, check_date = ? WHERE id = ?",
		item.Project, item.Catalog, item.ItemDesc, item.ItemCondition,item.ConnectionName, item.DashboardURL, item.Status, item.StatusDesc, item.Screen, item.CheckDate, item.ID)
	return err
}

// GetMonitorMetric retrieves a monitor item from the database
func GetMonitorMetric(id int64) (MonitorMetric, error) {
	var item MonitorMetric
	var checkDateBytes []byte
	err := DB.QueryRow("SELECT id, project, catalog, item_desc, item_condition,conn_name, dashboard_url, status, status_desc, screen, check_date FROM monitor_metrics WHERE id = ?", id).Scan(&item.ID, &item.Project, &item.Catalog, &item.ItemDesc, &item.ItemCondition,&item.ConnectionName, &item.DashboardURL, &item.Status, &item.StatusDesc, &item.Screen, &checkDateBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return MonitorMetric{}, fmt.Errorf("monitor item with ID %d not found", id)
		}
		return MonitorMetric{}, err
	}

	// Convert the byte slice to time.Time
	item.CheckDate, err = parseDateTime(string(checkDateBytes))
	if err != nil {
		return MonitorMetric{}, err
	}

	return item, nil
}

func GetMetricsByProject(project string) ([]MonitorMetric, error) {
	var items []MonitorMetric

	rows, err := DB.Query("SELECT id, project, catalog, item_desc, item_condition,conn_name, dashboard_url, status, status_desc, screen, check_date FROM monitor_metrics WHERE project=? ", project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item MonitorMetric
		var checkDateBytes []byte
		err = rows.Scan(&item.ID, &item.Project, &item.Catalog, &item.ItemDesc, &item.ItemCondition,&item.ConnectionName, &item.DashboardURL, &item.Status, &item.StatusDesc, &item.Screen, &checkDateBytes)
		if err != nil {
			return nil, err
		}

		item.CheckDate, err = parseDateTime(string(checkDateBytes))
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func GetMetrics(page int) ([]MonitorMetric, error) {
	var items []MonitorMetric
	offset := (page - 1) * itemsPerPage

	rows, err := DB.Query("SELECT id, project, catalog, item_desc, item_condition,conn_name, dashboard_url, status, status_desc, screen, check_date FROM monitor_metrics LIMIT ? OFFSET ?", itemsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item MonitorMetric
		var checkDateBytes []byte
		err = rows.Scan(&item.ID, &item.Project, &item.Catalog, &item.ItemDesc, &item.ItemCondition,&item.ConnectionName, &item.DashboardURL, &item.Status, &item.StatusDesc, &item.Screen, &checkDateBytes)
		if err != nil {
			return nil, err
		}

		item.CheckDate, err = parseDateTime(string(checkDateBytes))
		if err != nil {
			return nil, err
		}
		item.CheckDate = item.CheckDate.Local()
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// InsertMonitorJob inserts a new monitor config into the database
func InsertMonitorJob(config MonitorJob) (int64, error) {
	result, err := DB.Exec("INSERT INTO monitor_job (project, cron, enable) VALUES (?, ?, ?)",
		config.Project, config.Cron, config.Enable)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteMonitorJob deletes a monitor config from the database
func DeleteMonitorJob(id int64) error {
	_, err := DB.Exec("DELETE FROM monitor_job WHERE id = ?", id)
	return err
}

// UpdateMonitorJob updates a monitor config in the database
func UpdateMonitorJob(config MonitorJob) error {
	_, err := DB.Exec("UPDATE monitor_job SET project = ?, cron = ?, enable = ? WHERE id = ?",
		config.Project, config.Cron, config.Enable, config.ID)
	return err
}

// GetMonitorJob retrieves a monitor config from the database
func GetMonitorJob(id int64) (MonitorJob, error) {
	var config MonitorJob
	err := DB.QueryRow("SELECT id, project, cron, enable FROM monitor_job WHERE id = ?", id).Scan(&config.ID, &config.Project, &config.Cron, &config.Enable)
	if err != nil {
		if err == sql.ErrNoRows {
			return MonitorJob{}, fmt.Errorf("monitor config with ID %d not found", id)
		}
		return MonitorJob{}, err
	}

	return config, nil
}

// GetMonitorJobs retrieves all monitor configs from the database
func GetMonitorJobs() ([]MonitorJob, error) {
	var configs []MonitorJob

	rows, err := DB.Query("SELECT id, project, cron, enable FROM monitor_job")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var config MonitorJob
		err = rows.Scan(&config.ID, &config.Project, &config.Cron, &config.Enable)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func GetEnableMonitorJobs() ([]MonitorJob, error) {
	var configs []MonitorJob

	rows, err := DB.Query("SELECT id, project, cron, enable FROM monitor_job where enable=true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var config MonitorJob
		err = rows.Scan(&config.ID, &config.Project, &config.Cron, &config.Enable)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func GetDistinct(key, tableName string) ([]map[string]string, error) {
	var items []map[string]string

	rows, err := DB.Query(fmt.Sprintf("SELECT DISTINCT %s FROM %s", key,tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return nil, err
		}

		item := make(map[string]string)
		item[key] = value
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}


// InsertMonitorConnection inserts a new monitor connection into the database
func InsertMonitorConnection(conn MonitorConnection) (int64, error) {
    result, err := DB.Exec("INSERT INTO monitor_connections (conn_name, conn_username, conn_password, conn_url) VALUES (?, ?, ?, ?)",
        conn.Name, conn.Username, conn.Password, conn.URL)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return id, nil
}

// DeleteMonitorConnection deletes a monitor connection from the database
func DeleteMonitorConnection(id int64) error {
    _, err := DB.Exec("DELETE FROM monitor_connections WHERE id = ?", id)
    return err
}

// UpdateMonitorConnection updates a monitor connection in the database
func UpdateMonitorConnection(conn MonitorConnection) error {
    _, err := DB.Exec("UPDATE monitor_connections SET conn_name = ?, conn_username = ?, conn_password = ?, conn_url = ? WHERE id = ?",
        conn.Name, conn.Username, conn.Password, conn.URL, conn.ID)
    return err
}

// GetMonitorConnection retrieves a monitor connection from the database
func GetMonitorConnection(id int64) (MonitorConnection, error) {
    var conn MonitorConnection
    err := DB.QueryRow("SELECT id, conn_name, conn_username, conn_password, conn_url FROM monitor_connections WHERE id = ?", id).Scan(&conn.ID, &conn.Name, &conn.Username, &conn.Password, &conn.URL)
    if err != nil {
        if err == sql.ErrNoRows {
            return MonitorConnection{}, fmt.Errorf("monitor connection with ID %d not found", id)
        }
        return MonitorConnection{}, err
    }

    return conn, nil
}

func GetMonitorConnectionByName(connectionName string) (MonitorConnection, error) {
	fmt.Println(connectionName)
    var conn MonitorConnection
    err := DB.QueryRow("SELECT id, conn_name, conn_username, conn_password, conn_url FROM monitor_connections WHERE conn_name = ?", connectionName).Scan(&conn.ID, &conn.Name, &conn.Username, &conn.Password, &conn.URL)
    if err != nil {
        if err == sql.ErrNoRows {
            return MonitorConnection{}, fmt.Errorf("monitor connection with ID %d not found", connectionName)
        }
        return MonitorConnection{}, err
    }

    return conn, nil
}

// GetMonitorConnections retrieves all monitor connections from the database
func GetMonitorConnections() ([]MonitorConnection, error) {
    var conns []MonitorConnection

    rows, err := DB.Query("SELECT id, conn_name, conn_username, conn_password, conn_url FROM monitor_connections")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var conn MonitorConnection
        err = rows.Scan(&conn.ID, &conn.Name, &conn.Username, &conn.Password, &conn.URL)
        if err != nil {
            return nil, err
        }

        conns = append(conns, conn)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return conns, nil
}

//SELECT DISTINCT project FROM monitor_metrics;

