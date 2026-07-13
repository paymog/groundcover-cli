package raw

var StorageManagementCommands = []Command{
	{
		Name:        []string{"storage-management", "get"},
		Method:      "GET",
		Path:        "/api/storage-management/:dataType",
		Description: "Get retention and indexing settings for a data type",
		PathParams:  []string{"dataType"},
	},
	{
		Name:        []string{"storage-management", "update"},
		Method:      "PUT",
		Path:        "/api/storage-management/:dataType",
		Description: "Update retention and indexing settings for a data type",
		PathParams:  []string{"dataType"},
	},
}
