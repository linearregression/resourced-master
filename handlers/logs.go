package handlers

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/jmoiron/sqlx"

	"github.com/resourced/resourced-master/dal"
	"github.com/resourced/resourced-master/libhttp"
)

func GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	qParams := r.URL.Query()

	toString := qParams.Get("To")
	if toString == "" {
		toString = qParams.Get("to")
	}
	to, err := strconv.ParseInt(toString, 10, 64)

	fromString := qParams.Get("From")
	if fromString == "" {
		fromString = qParams.Get("from")
	}
	from, err := strconv.ParseInt(fromString, 10, 64)

	query := qParams.Get("q")

	currentUser := context.Get(r, "currentUser").(*dal.UserRow)

	currentCluster := context.Get(r, "currentCluster").(*dal.ClusterRow)

	db := context.Get(r, "db.Core").(*sqlx.DB)

	tsLogsDB := context.Get(r, "db.TSLog").(*sqlx.DB)

	// -----------------------------------
	// Create channels to receive SQL rows
	// -----------------------------------
	tsLogsChan := make(chan *dal.TSLogRowsWithError)
	defer close(tsLogsChan)

	savedQueriesChan := make(chan *dal.SavedQueryRowsWithError)
	defer close(savedQueriesChan)

	// --------------------------
	// Fetch SQL rows in parallel
	// --------------------------
	go func(currentCluster *dal.ClusterRow, query string) {
		var tsLogs []*dal.TSLogRow
		var err error

		if fromString == "" && toString == "" {
			tsLogs, err = dal.NewTSLog(tsLogsDB).AllByClusterIDLastRowIntervalAndQuery(nil, currentCluster.ID, "15 minute", query)

			if len(tsLogs) > 0 {
				if from == 0 {
					from = tsLogs[0].Created.UTC().Unix()
				}
				if to == 0 {
					to = tsLogs[len(tsLogs)-1].Created.UTC().Unix()
				}
			}

		} else {
			tsLogs, err = dal.NewTSLog(tsLogsDB).AllByClusterIDRangeAndQuery(nil, currentCluster.ID, from, to, query)
		}

		tsLogRowsWithError := &dal.TSLogRowsWithError{}
		tsLogRowsWithError.TSLogRows = tsLogs
		tsLogRowsWithError.Error = err
		tsLogsChan <- tsLogRowsWithError
	}(currentCluster, query)

	go func(currentCluster *dal.ClusterRow) {
		savedQueriesWithError := &dal.SavedQueryRowsWithError{}
		savedQueriesWithError.SavedQueries, savedQueriesWithError.Error = dal.NewSavedQuery(db).AllByClusterIDAndType(nil, currentCluster.ID, "hosts")
		savedQueriesChan <- savedQueriesWithError
	}(currentCluster)

	// -----------------------------------
	// Wait for channels to return results
	// -----------------------------------
	tsLogsWithError := <-tsLogsChan
	if tsLogsWithError.Error != nil && tsLogsWithError.Error.Error() != "sql: no rows in result set" {
		libhttp.HandleErrorJson(w, tsLogsWithError.Error)
		return
	}

	savedQueriesWithError := <-savedQueriesChan
	if savedQueriesWithError.Error != nil && savedQueriesWithError.Error.Error() != "sql: no rows in result set" {
		libhttp.HandleErrorJson(w, savedQueriesWithError.Error)
		return
	}

	data := struct {
		CSRFToken          string
		Addr               string
		CurrentUser        *dal.UserRow
		Clusters           []*dal.ClusterRow
		CurrentClusterJson string
		Logs               []*dal.TSLogRow
		SavedQueries       []*dal.SavedQueryRow
		From               int64
		To                 int64
	}{
		csrf.Token(r),
		context.Get(r, "addr").(string),
		currentUser,
		context.Get(r, "clusters").([]*dal.ClusterRow),
		string(context.Get(r, "currentClusterJson").([]byte)),
		tsLogsWithError.TSLogRows,
		savedQueriesWithError.SavedQueries,
		from,
		to,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/logs/list.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}

func GetLogsExecutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	qParams := r.URL.Query()

	toString := qParams.Get("To")
	if toString == "" {
		toString = qParams.Get("to")
	}
	to, err := strconv.ParseInt(toString, 10, 64)
	if err != nil {
		to = time.Now().UTC().Unix()
	}

	fromString := qParams.Get("From")
	if fromString == "" {
		fromString = qParams.Get("from")
	}
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		// default is 15 minutes range
		from = to - 900
	}

	query := qParams.Get("q")

	currentUser := context.Get(r, "currentUser").(*dal.UserRow)

	currentCluster := context.Get(r, "currentCluster").(*dal.ClusterRow)

	db := context.Get(r, "db.Core").(*sqlx.DB)

	tsExecutorLogsDB := context.Get(r, "db.TSExecutorLog").(*sqlx.DB)

	// -----------------------------------
	// Create channels to receive SQL rows
	// -----------------------------------
	tsLogsChan := make(chan *dal.TSExecutorLogRowsWithError)
	defer close(tsLogsChan)

	savedQueriesChan := make(chan *dal.SavedQueryRowsWithError)
	defer close(savedQueriesChan)

	// --------------------------
	// Fetch SQL rows in parallel
	// --------------------------
	go func(currentCluster *dal.ClusterRow, from, to int64, query string) {
		tsLogRowsWithError := &dal.TSExecutorLogRowsWithError{}
		tsLogRowsWithError.TSExecutorLogRows, tsLogRowsWithError.Error = dal.NewTSExecutorLog(tsExecutorLogsDB).AllByClusterIDRangeAndQuery(nil, currentCluster.ID, from, to, query)
		tsLogRowsWithError.Error = err
		tsLogsChan <- tsLogRowsWithError
	}(currentCluster, from, to, query)

	go func(currentCluster *dal.ClusterRow) {
		savedQueriesWithError := &dal.SavedQueryRowsWithError{}
		savedQueriesWithError.SavedQueries, savedQueriesWithError.Error = dal.NewSavedQuery(db).AllByClusterIDAndType(nil, currentCluster.ID, "hosts")
		savedQueriesChan <- savedQueriesWithError
	}(currentCluster)

	// -----------------------------------
	// Wait for channels to return results
	// -----------------------------------
	tsLogsWithError := <-tsLogsChan
	if tsLogsWithError.Error != nil && tsLogsWithError.Error.Error() != "sql: no rows in result set" {
		libhttp.HandleErrorJson(w, tsLogsWithError.Error)
		return
	}

	savedQueriesWithError := <-savedQueriesChan
	if savedQueriesWithError.Error != nil && savedQueriesWithError.Error.Error() != "sql: no rows in result set" {
		libhttp.HandleErrorJson(w, savedQueriesWithError.Error)
		return
	}

	data := struct {
		CSRFToken          string
		Addr               string
		CurrentUser        *dal.UserRow
		Clusters           []*dal.ClusterRow
		CurrentClusterJson string
		Logs               []*dal.TSExecutorLogRow
		SavedQueries       []*dal.SavedQueryRow
		From               int64
		To                 int64
	}{
		csrf.Token(r),
		context.Get(r, "addr").(string),
		currentUser,
		context.Get(r, "clusters").([]*dal.ClusterRow),
		string(context.Get(r, "currentClusterJson").([]byte)),
		tsLogsWithError.TSExecutorLogRows,
		savedQueriesWithError.SavedQueries,
		from,
		to,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/logs/list.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}

func PostApiLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accessTokenRow := context.Get(r, "accessTokenRow").(*dal.AccessTokenRow)

	dataJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tsLogsDB := context.Get(r, "db.TSLog").(*sqlx.DB)

	err = dal.NewTSLog(tsLogsDB).CreateFromJSON(nil, accessTokenRow.ClusterID, dataJson)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write([]byte(`{"Message": "Success"}`))
}
