package dal

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	sqlx_types "github.com/jmoiron/sqlx/types"
)

func NewWatcher(db *sqlx.DB) *Watcher {
	watcher := &Watcher{}
	watcher.db = db
	watcher.table = "watchers"
	watcher.hasID = true

	return watcher
}

type WatcherRowsWithError struct {
	Watchers []*WatcherRow
	Error    error
}

type WatcherRow struct {
	BaseRow
	ID               int64               `db:"id" json:"-"`
	ClusterID        int64               `db:"cluster_id"`
	SavedQuery       sql.NullString      `db:"saved_query"`
	Name             string              `db:"name"`
	LowAffectedHosts int64               `db:"low_affected_hosts"`
	HostsLastUpdated string              `db:"hosts_last_updated"`
	CheckInterval    string              `db:"check_interval"`
	IsSilenced       bool                `db:"is_silenced"`
	ActiveCheck      sqlx_types.JSONText `db:"active_check"`
}

func (wr *WatcherRow) IsPassive() bool {
	return wr.SavedQuery.String != ""
}

func (wr *WatcherRow) HostsLastUpdatedForPostgres() string {
	return strings.Replace(wr.HostsLastUpdated, " ago", "", -1)
}

func (wr *WatcherRow) Command() string {
	return wr.JSONAttrString(wr.ActiveCheck, "Command")
}

func (wr *WatcherRow) SSHUser() string {
	return wr.JSONAttrString(wr.ActiveCheck, "SSHUser")
}

func (wr *WatcherRow) SSHPort() string {
	return wr.JSONAttrString(wr.ActiveCheck, "SSHPort")
}

func (wr *WatcherRow) HTTPHeadersString() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPHeaders")
}

func (wr *WatcherRow) HTTPHeaders() map[string]string {
	data := make(map[string]string)

	asString := wr.HTTPHeadersString()
	asList := strings.Split(asString, ",")

	for _, kvString := range asList {
		kvList := strings.Split(kvString, ":")
		if len(kvList) >= 2 {
			data[strings.TrimSpace(kvList[0])] = strings.TrimSpace(kvList[1])
		}
	}

	return data
}

func (wr *WatcherRow) HTTPPostBody() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPPostBody")
}

func (wr *WatcherRow) HTTPMethod() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPMethod")
}

func (wr *WatcherRow) HTTPScheme() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPScheme")
}

func (wr *WatcherRow) HTTPPort() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPPort")
}

func (wr *WatcherRow) HTTPPath() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPPath")
}

func (wr *WatcherRow) HTTPCode() int {
	code := wr.JSONAttrFloat64(wr.ActiveCheck, "HTTPCode")
	return int(code)
}

func (wr *WatcherRow) HTTPUser() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HTTPUser")
}

func (wr *WatcherRow) HTTPPass() string {
	inBase64 := wr.JSONAttrString(wr.ActiveCheck, "HTTPPass")
	decryptedBytes, err := base64.StdEncoding.DecodeString(inBase64)
	if err != nil {
		return ""
	}
	return string(decryptedBytes)
}

func (wr *WatcherRow) HostsListString() string {
	return wr.JSONAttrString(wr.ActiveCheck, "HostsList")
}

func (wr *WatcherRow) HostsList() []string {
	if wr.HostsListString() == "" {
		return make([]string, 0)
	}

	return strings.Split(wr.HostsListString(), "\n")
}

func (wr *WatcherRow) PerformActiveCheckPing(hostname string) (outBytes []byte, err error) {
	return exec.Command("ping", "-c", "1", hostname).CombinedOutput()
}

func (wr *WatcherRow) PerformActiveCheckSSH(hostname string) (outBytes []byte, err error) {
	sshOptions := []string{"-o BatchMode=yes", "-o ConnectTimeout=10"}

	if wr.SSHPort() != "" {
		sshOptions = append(sshOptions, []string{"-p", wr.SSHPort()}...)
	}

	userAtHost := hostname

	if wr.SSHUser() != "" {
		userAtHost = fmt.Sprintf("%v@%v", wr.SSHUser(), hostname)
	}

	sshOptions = append(sshOptions, userAtHost)

	return exec.Command("ssh", sshOptions...).CombinedOutput()
}

func (wr *WatcherRow) PerformActiveCheckHTTP(hostname string) (resp *http.Response, err error) {
	url := fmt.Sprintf("%v://%v:%s", wr.HTTPScheme(), hostname, wr.HTTPPort())

	client := &http.Client{}

	var req *http.Request

	if wr.HTTPPostBody() != "" {
		req, err = http.NewRequest(strings.ToUpper(wr.HTTPMethod()), url, strings.NewReader(wr.HTTPPostBody()))

		// Detect if POST body is JSON and set content-type
		if strings.HasPrefix(wr.HTTPPostBody(), "{") || strings.HasPrefix(wr.HTTPPostBody(), "[") {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

	} else {
		req, err = http.NewRequest(strings.ToUpper(wr.HTTPMethod()), url, nil)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method":     "http.NewRequest",
			"URL":        url,
			"HTTPMethod": wr.HTTPMethod(),
		}).Error(err)
		return nil, err
	}

	for headerKey, headerVal := range wr.HTTPHeaders() {
		req.Header.Add(headerKey, headerVal)
	}

	if wr.HTTPUser() != "" || wr.HTTPPass() != "" {
		req.SetBasicAuth(wr.HTTPUser(), wr.HTTPPass())
	}

	resp, err = client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Method":     "http.Client{}.Do",
			"URL":        url,
			"HTTPMethod": wr.HTTPMethod(),
		}).Error(err)
		return nil, err

	} else if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	return resp, err
}

type Watcher struct {
	Base
}

// All returns all watchers rows.
func (w *Watcher) All(tx *sqlx.Tx) ([]*WatcherRow, error) {
	rows := []*WatcherRow{}
	query := fmt.Sprintf("SELECT * FROM %v", w.table)
	err := w.db.Select(&rows, query)

	return rows, err
}

// AllSplitToDaemons returns all watchers rows divided into daemons equally.
func (w *Watcher) AllSplitToDaemons(tx *sqlx.Tx, daemons []string) (map[string][]*WatcherRow, error) {
	wrs, err := w.All(tx)
	if err != nil {
		return nil, err
	}

	buckets := make([][]*WatcherRow, len(daemons))
	for i, _ := range daemons {
		buckets[i] = make([]*WatcherRow, 0)
	}

	bucketsPointer := 0
	for _, wr := range wrs {
		buckets[bucketsPointer] = append(buckets[bucketsPointer], wr)
		bucketsPointer = bucketsPointer + 1

		if bucketsPointer >= len(buckets) {
			bucketsPointer = 0
		}
	}

	result := make(map[string][]*WatcherRow)

	for i, watchersInbucket := range buckets {
		result[daemons[i]] = watchersInbucket
	}

	return result, err
}

// AllPassiveByClusterID returns all rows by cluster_id.
func (w *Watcher) AllPassiveByClusterID(tx *sqlx.Tx, clusterID int64) ([]*WatcherRow, error) {
	rows := []*WatcherRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE cluster_id=$1 AND saved_query <> '' ORDER BY name ASC", w.table)
	err := w.db.Select(&rows, query, clusterID)

	return rows, err
}

// AllActiveByClusterID returns all rows by cluster_id and command != nil.
func (w *Watcher) AllActiveByClusterID(tx *sqlx.Tx, clusterID int64) ([]*WatcherRow, error) {
	rows := []*WatcherRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE cluster_id=$1 AND saved_query = '' ORDER BY name ASC", w.table)
	err := w.db.Select(&rows, query, clusterID)

	return rows, err
}

func (w *Watcher) rowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*WatcherRow, error) {
	id, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return w.GetByID(tx, id)
}

// GetByID returns one record by id.
func (w *Watcher) GetByID(tx *sqlx.Tx, id int64) (*WatcherRow, error) {
	wr := &WatcherRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", w.table)
	err := w.db.Get(wr, query, id)

	return wr, err
}

// CreateOrUpdateParameters builds params for insert or update.
func (w *Watcher) CreateOrUpdateParameters(clusterID int64, savedQuery, name string, lowAffectedHosts int64, hostsLastUpdated, checkInterval string, activeCheck map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	data["cluster_id"] = clusterID
	data["saved_query"] = savedQuery
	data["name"] = name
	data["low_affected_hosts"] = lowAffectedHosts
	data["hosts_last_updated"] = hostsLastUpdated
	data["check_interval"] = checkInterval

	activeCheckJsonBytes, err := json.Marshal(activeCheck)
	if err == nil {
		data["active_check"] = activeCheckJsonBytes
	}

	return data
}

// Create inserts one row
func (w *Watcher) Create(tx *sqlx.Tx, data map[string]interface{}) (*WatcherRow, error) {
	sqlResult, err := w.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return w.rowFromSqlResult(tx, sqlResult)
}

// HostsToPerformActiveChecks fetches list of hostnames to perform active checks on
func (w *Watcher) HostsToPerformActiveChecks(clusterID int64, watcherRow *WatcherRow) (hostnames []string) {
	// Checking hosts based on external lists.
	if len(watcherRow.HostsList()) > 0 {
		hostnames = watcherRow.HostsList()
	} else {
		hosts, err := NewHost(w.db).AllByClusterIDAndUpdatedInterval(nil, clusterID, watcherRow.HostsLastUpdatedForPostgres())
		if err != nil {
			return make([]string, 0)
		}

		hostnames = make([]string, len(hosts))

		for i, host := range hosts {
			hostnames[i] = host.Hostname
		}
	}

	return hostnames
}
