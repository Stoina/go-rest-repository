package restrepo

// SQLRepositorySetting exported
// ...
type SQLRepositorySetting struct {
	TableName 			string
	IDColumn 			string
	ParameterMapping 	map[string]string
}

// NewSQLRepositorySetting exported
// ...
func NewSQLRepositorySetting(tableName string, idColumn string, parameterMapping map[string]string) *SQLRepositorySetting {
	return &SQLRepositorySetting{
		TableName: tableName,
		IDColumn: idColumn,
		ParameterMapping: parameterMapping};
}