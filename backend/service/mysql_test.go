package service

import (
	"testing"
)

func TestConnectDB(t *testing.T) {
	// Test the ConnectDB function
	InitDBConnection()
	defer DB.Close()
}

func TestInsertMonitorItem(t *testing.T) {
	// Test the InsertMonitorItem function
	item := MonitorMetric{
		Project:       "Test Project",
		Catalog:       "Test Catalog",
		ItemDesc:      "Test Description",
		ItemCondition: "Test Condition",
		DashboardURL:  "https://example.com/test",
		Status:        true,
		StatusDesc:    "",
		Screen:        "",
	}

	InitDBConnection()

	id, err := InsertMonitorMetric(item)
	if err != nil {
		t.Errorf("Failed to insert monitor item: %v", err)
	}
	t.Log(id)
	// Verify the inserted item
	retrievedItem, err := GetMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to retrieve monitor item: %v", err)
	}

	if retrievedItem.ItemDesc != item.ItemDesc {
		t.Errorf("Retrieved monitor item does not match the inserted item")
	}

}

func TestUpdateMonitorItem(t *testing.T) {
	// Test the UpdateMonitorItem function
	item := MonitorMetric{
		Project:       "Test Project",
		Catalog:       "Test Catalog",
		ItemDesc:      "Test Description",
		ItemCondition: "Test Condition",
		ConnectionName: "default_grafana",
		DashboardURL:  "https://example.com/test",
		Status:        true,
		StatusDesc:    "OK",
		Screen:        "test_screen.png",
	}

	InitDBConnection()

	id, err := InsertMonitorMetric(item)
	if err != nil {
		t.Errorf("Failed to insert monitor item: %v", err)
	}

	// Update the item
	item.StatusDesc = "Updated"
	item.ID = id

	err = UpdateMonitorMetric(item)
	if err != nil {
		t.Errorf("Failed to update monitor item: %v", err)
	}

	// Verify the updated item
	t.Log(id)
	updatedItem, err := GetMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to retrieve monitor item: %v", err)
	}

	if updatedItem.StatusDesc != "Updated" {
		t.Errorf("Monitor item was not updated correctly")
	}

	// // Clean up by deleting the inserted item
	err = DeleteMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to delete monitor item: %v", err)
	}
}

func TestDeleteMonitorItem(t *testing.T) {
	// Test the DeleteMonitorItem function
	item := MonitorMetric{
		Project:       "Test Project",
		Catalog:       "Test Catalog",
		ItemDesc:      "Test Description",
		ItemCondition: "Test Condition",
		DashboardURL:  "https://example.com/test",
		Status:        true,
		StatusDesc:    "OK",
		Screen:        "test_screen.png",
	}

	InitDBConnection()

	id, err := InsertMonitorMetric(item)
	if err != nil {
		t.Errorf("Failed to insert monitor item: %v", err)
	}

	// Delete the item
	err = DeleteMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to delete monitor item: %v", err)
	}

	// Verify that the item is deleted
	_, err = GetMonitorMetric(id)
	if err == nil {
		t.Errorf("Monitor item was not deleted correctly")
	}
}

func TestGetMonitorItem(t *testing.T) {
	// Test the GetMonitorItem function

	item := MonitorMetric{
		Project:       "Test Project",
		Catalog:       "Test Catalog",
		ItemDesc:      "Test Description",
		ItemCondition: "Test Condition",
		DashboardURL:  "https://example.com/test",
		Status:        true,
		StatusDesc:    "OK",
		Screen:        "test_screen.png",
	}

	InitDBConnection()
	id, err := InsertMonitorMetric(item)
	if err != nil {
		t.Errorf("Failed to insert monitor item: %v", err)
	}

	// Retrieve the item
	_, err = GetMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to retrieve monitor item: %v", err)
	}

	// Test for non-existent ID
	_, err = GetMonitorMetric(999999)
	if err == nil {
		t.Errorf("GetMonitorItem should return an error for non-existent ID")
	}

	// Clean up by deleting the inserted item
	err = DeleteMonitorMetric(id)
	if err != nil {
		t.Errorf("Failed to delete monitor item: %v", err)
	}
}
