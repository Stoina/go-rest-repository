package restrepo

import (
	"encoding/json"
	"net/url"
	"strings"

	db "github.com/Stoina/go-database"
)

// SQLRepository exported
// SQLRepository ...
type SQLRepository struct {
	name         	string
	url          	string
	dbConnection 	*db.Connection
	setting 		*SQLRepositorySetting
}

// NewSQLRepository exported
// NewSQLRepository ...
func NewSQLRepository(name string, url string, dbConn *db.Connection, setting *SQLRepositorySetting) *SQLRepository {
	return &SQLRepository{
		name:         name,
		url:          url,
		dbConnection: dbConn,
		setting: 	  setting}
}

// Post exported
// With a post request to the container resource you can create a new resource.
func (sqlRepo SQLRepository) Post(contentType string, content string) *RepositoryResult {
	if contentType == "application/json" {
		return insertJSON(sqlRepo, content)
	}

	return nil
}

// Put exported
// With a put request to the container resource you can overwrite the resource with the representation in the request.
func (sqlRepo SQLRepository) Put(par string) *RepositoryResult {
	return nil
}

// Patch exported
// With the http patch method, individual properties of a resource can be manipulated in a targeted manner.
func (sqlRepo SQLRepository) Patch(par string) *RepositoryResult {
	return nil
}

// Delete exported
// With this method an existing resource can be deleted
func (sqlRepo SQLRepository) Delete(par string) *RepositoryResult {
	return nil
}

// Get exported
// Get ...
func (sqlRepo SQLRepository) Get(calledURL *url.URL) *RepositoryResult {
	query := buildQueryString(&sqlRepo, calledURL)
	return executeQuery(sqlRepo.dbConnection, query)
}

// Name exported
// Name ....
func (sqlRepo SQLRepository) Name() string {
	return sqlRepo.name
}

// URL exported
// URL ...
func (sqlRepo SQLRepository) URL() string {
	return sqlRepo.url
}

func buildQueryString(sqlRepo *SQLRepository, calledURL *url.URL) string {
	query := "select * from \"" + sqlRepo.setting.TableName + "\""
	
	urlParameter := calledURL.Query()
	if (len(urlParameter) > 0) {
		queryConditions := getQueryParameterFromURL(sqlRepo.setting, urlParameter)

		for i, queryCondition := range queryConditions {
			condition := " and "

			if i == 0 {
				condition = " where "
			}

			query += condition + queryCondition
		}

	} else {
		urlID := getIDParameterFromURL(sqlRepo.URL(), calledURL)
		
		if urlID != "" {
			query += " where \"" + sqlRepo.setting.IDColumn + "\" = " + urlID
		}
	}

	return query
}

func getIDParameterFromURL(ownURL string, calledURL *url.URL) string {
	splittedURL := strings.Split(calledURL.RequestURI(), "/")

	for _, value := range splittedURL {
		if value != ownURL && value != "" {
			return value
		} 
	}

	return ""
}

func getQueryParameterFromURL(sqlRepoSetting *SQLRepositorySetting, mappedURLValues url.Values) []string {

	queryConditions := make([]string, len(mappedURLValues))

	i := 0
	for urlValueKey, urlValues := range mappedURLValues {
		sqlColumnName := sqlRepoSetting.ParameterMapping[urlValueKey]

		if sqlColumnName != "" {
			conditionValue := ""

			for _, urlValue := range urlValues {
				conditionValue += urlValue
			}

			queryConditions[i] = "\"" + sqlColumnName + "\" = '" + conditionValue + "'"
		}

		i++
	}

	return queryConditions
}

func executeQuery(dbConn *db.Connection, query string) *RepositoryResult {

	var resultError error
	responseMessage := ""
	responseData := ""

	queryResult, resultError := dbConn.Query(query)

	if resultError == nil {
		data, resultError := queryResult.ConvertToJSON()
		
		if resultError == nil {
			responseMessage = "Data loaded successfully"
			responseData = data
		}	
	}
	
	if resultError == nil {
		return NewRepositoryResult(responseData, false, "", responseMessage, true)
	} 
	
	return NewRepositoryResult(responseData, true, resultError.Error(), responseMessage, false)
}

func insertJSON(sqlRepository SQLRepository, jsonContent string) *RepositoryResult {

	var jsonValues map[string]interface{}
	json.Unmarshal([]byte(jsonContent), &jsonValues)
	
	columndAndValueCount := len(jsonValues)

	columns := make([]string, columndAndValueCount) 
	values := make([]interface{}, columndAndValueCount)

	index := 0
	for key, value := range jsonValues {
		columns[index] = key
		values[index] = value

		index++
	}
	
	var resultError error
	responseMessage := ""
	responseData := ""

	insertStatement := db.NewInsertStatement(sqlRepository.setting.TableName, columns, values)
	dbResult, resultError := sqlRepository.dbConnection.Insert(insertStatement)

	if resultError == nil {
		responseData, resultError = dbResult.ConvertToJSON()
	}
	
	if resultError == nil {
		return NewRepositoryResult(responseData, false, "", responseMessage, true)
	} 
	
	return NewRepositoryResult(responseData, true, resultError.Error(), responseMessage, false)
}