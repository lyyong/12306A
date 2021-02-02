package init_data

import "testing"

func TestCreateAllTables(t *testing.T) {
	CreateAllTables()
}

func TestDropAllTables(t *testing.T) {
	DropAllTables()
}

func TestCreateTableTicketPool(t *testing.T) {
	CreateTableTicketPool()
}
