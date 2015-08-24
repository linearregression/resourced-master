package dal

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewSavedQuery(db *sqlx.DB) *SavedQuery {
	savedQuery := &SavedQuery{}
	savedQuery.db = db
	savedQuery.table = "saved_queries"
	savedQuery.hasID = true

	return savedQuery
}

type SavedQueryRow struct {
	ID        int64  `db:"id"`
	UserID    int64  `db:"user_id"`
	ClusterID int64  `db:"cluster_id"`
	Query     string `db:"query"`
}

type SavedQuery struct {
	Base
}

func (sq *SavedQuery) savedQueryRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*SavedQueryRow, error) {
	savedQueryId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return sq.GetByID(tx, savedQueryId)
}

// DeleteByID deletes record by id.
func (sq *SavedQuery) DeleteByID(tx *sqlx.Tx, id int64) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id=$1", sq.table)
	_, err := sq.db.Exec(query, id)

	return err
}

// AllByAccessToken returns all saved_query rows.
func (sq *SavedQuery) AllByAccessToken(tx *sqlx.Tx, accessTokenRow *AccessTokenRow) ([]*SavedQueryRow, error) {
	savedQueries := []*SavedQueryRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE cluster_id=$1", sq.table)
	err := sq.db.Select(&savedQueries, query, accessTokenRow.ClusterID)

	return savedQueries, err
}

// GetByID returns record by id.
func (sq *SavedQuery) GetByID(tx *sqlx.Tx, id int64) (*SavedQueryRow, error) {
	savedQueryRow := &SavedQueryRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", sq.table)
	err := sq.db.Get(savedQueryRow, query, id)

	return savedQueryRow, err
}

// GetByAccessTokenAndQuery returns record by savedQuery.
func (sq *SavedQuery) GetByAccessTokenAndQuery(tx *sqlx.Tx, accessTokenRow *AccessTokenRow, savedQuery string) (*SavedQueryRow, error) {
	savedQueryRow := &SavedQueryRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE user_id=$1 AND cluster_id=$2 AND query=$3", sq.table)
	err := sq.db.Get(savedQueryRow, query, accessTokenRow.UserID, accessTokenRow.ClusterID, savedQuery)

	return savedQueryRow, err
}

// CreateOrUpdate performs insert/update for one savedQuery data.
func (sq *SavedQuery) CreateOrUpdate(tx *sqlx.Tx, accessTokenRow *AccessTokenRow, savedQuery string) (*SavedQueryRow, error) {
	savedQueryRow, err := sq.GetByAccessTokenAndQuery(tx, accessTokenRow, savedQuery)

	data := make(map[string]interface{})
	data["user_id"] = accessTokenRow.UserID
	data["cluster_id"] = accessTokenRow.ClusterID
	data["query"] = savedQuery

	// Perform INSERT
	if savedQueryRow == nil || err != nil {
		sqlResult, err := sq.InsertIntoTable(tx, data)
		if err != nil {
			return nil, err
		}

		return sq.savedQueryRowFromSqlResult(tx, sqlResult)
	}

	return savedQueryRow, nil
}
