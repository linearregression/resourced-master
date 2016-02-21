package handlers

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"

	"github.com/resourced/resourced-master/dal"
	"github.com/resourced/resourced-master/libhttp"
	"github.com/resourced/resourced-master/multidb"
)

func PostMetrics(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db.Core").(*sqlx.DB)

	vars := mux.Vars(r)

	clusterIDString := vars["clusterid"]
	clusterID, err := strconv.ParseInt(clusterIDString, 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	key := r.FormValue("Key")

	_, err = dal.NewMetric(db).CreateOrUpdate(nil, clusterID, key)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/", 301)
}

func GetApiTSMetricsByHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := context.Get(r, "db.Core").(*sqlx.DB)

	id, err := getInt64SlugFromPath(w, r, "id")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	qParams := r.URL.Query()

	createdInterval := qParams.Get("CreatedInterval")
	if createdInterval == "" {
		createdInterval = "1 hour"
	}

	fromString := qParams.Get("From")
	if fromString == "" {
		fromString = qParams.Get("from")
	}
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		from = -1
	}

	toString := qParams.Get("To")
	if toString == "" {
		toString = qParams.Get("to")
	}
	to, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		to = -1
	}

	host := mux.Vars(r)["host"]

	metricRow, err := dal.NewMetric(db).GetById(nil, id)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tsMetricsDB := context.Get(r, "multidb.TSMetrics").(*multidb.MultiDB).PickRandom()

	var hcMetrics *dal.TSMetricHighchartPayload

	if from > 0 && to > 0 {
		hcMetrics, err = dal.NewTSMetric(tsMetricsDB).AllByMetricIDHostAndRangeForHighchart(nil, metricRow.ClusterID, id, host, from, to)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

	} else {
		hcMetrics, err = dal.NewTSMetric(tsMetricsDB).AllByMetricIDHostAndIntervalForHighchart(nil, metricRow.ClusterID, id, host, createdInterval)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
	}

	hcMetricsJSON, err := json.Marshal(hcMetrics)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(hcMetricsJSON)
}

func GetApiTSMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := context.Get(r, "db.Core").(*sqlx.DB)

	qParams := r.URL.Query()

	createdInterval := qParams.Get("CreatedInterval")
	if createdInterval == "" {
		createdInterval = "1 hour"
	}

	fromString := qParams.Get("From")
	if fromString == "" {
		fromString = qParams.Get("from")
	}
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		from = -1
	}

	toString := qParams.Get("To")
	if toString == "" {
		toString = qParams.Get("to")
	}
	to, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		to = -1
	}

	id, err := getInt64SlugFromPath(w, r, "id")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	metricRow, err := dal.NewMetric(db).GetById(nil, id)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tsMetricsDB := context.Get(r, "multidb.TSMetrics").(*multidb.MultiDB).PickRandom()

	var hcMetrics []*dal.TSMetricHighchartPayload

	if from > 0 && to > 0 {
		hcMetrics, err = dal.NewTSMetric(tsMetricsDB).AllByMetricIDAndRangeForHighchart(nil, metricRow.ClusterID, id, from, to)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

	} else {
		hcMetrics, err = dal.NewTSMetric(tsMetricsDB).AllByMetricIDAndIntervalForHighchart(nil, metricRow.ClusterID, id, createdInterval)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
	}

	hcMetricsJSON, err := json.Marshal(hcMetrics)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(hcMetricsJSON)
}

func GetApiTSMetrics15Min(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := context.Get(r, "db.Core").(*sqlx.DB)

	qParams := r.URL.Query()

	createdInterval := qParams.Get("CreatedInterval")
	if createdInterval == "" {
		createdInterval = "1 hour"
	}

	fromString := qParams.Get("From")
	if fromString == "" {
		fromString = qParams.Get("from")
	}
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		from = -1
	}

	toString := qParams.Get("To")
	if toString == "" {
		toString = qParams.Get("to")
	}
	to, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		to = -1
	}

	id, err := getInt64SlugFromPath(w, r, "id")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	metricRow, err := dal.NewMetric(db).GetById(nil, id)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tsMetricsDB := context.Get(r, "multidb.TSMetrics").(*multidb.MultiDB).PickRandom()

	var hcMetrics []*dal.TSMetricHighchartPayload

	if from > 0 && to > 0 {
		hcMetrics, err = dal.NewTSMetricAggr15m(tsMetricsDB).AllByMetricIDAndRangeForHighchart(nil, metricRow.ClusterID, id, from, to)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

	} else {
		hcMetrics, err = dal.NewTSMetricAggr15m(tsMetricsDB).AllByMetricIDAndIntervalForHighchart(nil, metricRow.ClusterID, id, createdInterval)
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}
	}

	hcMetricsJSON, err := json.Marshal(hcMetrics)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(hcMetricsJSON)
}

// PostPutDeleteMetricID handles POST, PUT, and DELETE
func PostPutDeleteMetricID(w http.ResponseWriter, r *http.Request) {
	method := r.FormValue("_method")
	if method == "" {
		method = "put"
	}

	if method == "post" || method == "put" {
		PutMetricID(w, r)
	} else if method == "delete" {
		DeleteMetricID(w, r)
	}
}

// PutMetricID is not supported
func PutMetricID(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 301)
}

// DeleteMetricID deletes metrics by ID
func DeleteMetricID(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db.Core").(*sqlx.DB)

	vars := mux.Vars(r)

	clusterIDString := vars["clusterid"]
	clusterID, err := strconv.ParseInt(clusterIDString, 10, 64)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	id, err := getInt64SlugFromPath(w, r, "id")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	_, err = dal.NewMetric(db).DeleteByClusterIDAndID(nil, clusterID, id)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = dal.NewGraph(db).DeleteMetricFromGraphs(nil, clusterID, id)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/", 301)
}
