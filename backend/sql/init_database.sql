CREATE DATABASE bedrock_claude3_grafana_automonitor;
USE bedrock_claude3_grafana_automonitor;

CREATE TABLE  IF NOT EXISTS  monitor_metrics (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `project` VARCHAR(255) NOT NULL,
    `catalog` VARCHAR(255) NOT NULL,
    `item_desc` TEXT NOT NULL,
    `item_condition` TEXT NOT NULL,
    `conn_name` varchar(255) DEFAULT NULL,
    `dashboard_url` VARCHAR(255) NOT NULL,
    `status` BOOLEAN NOT NULL,
    `status_desc` TEXT NOT NULL,
    `screen` VARCHAR(255) NOT NULL,
    `check_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) AUTO_INCREMENT = 100000;

CREATE TABLE   IF NOT EXISTS  monitor_job (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project VARCHAR(255) NOT NULL,
    cron VARCHAR(255) NOT NULL,
    enable BOOLEAN NOT NULL
) AUTO_INCREMENT = 100000;

CREATE TABLE  IF NOT EXISTS  `monitor_connections` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conn_name` varchar(255) NOT NULL UNIQUE,
  `conn_username` varchar(255) NOT NULL,
  `conn_password` varchar(255) NOT NULL,
  `conn_url` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=10000;

INSERT INTO monitor_connections (conn_name, conn_username, conn_password, conn_url)
VALUES ('default_grafana', 'admin', 'your-secret-pw', 'http://localhost:3100');