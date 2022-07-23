package isuports

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"

	// "github.com/mattn/go-sqlite3"
	// proxy "github.com/shogo82148/go-sql-proxy"
) 

//テナントDBのパスを返す
func tenantDBPath(id int64) string {
	tenantDBDir := getEnv("ISUCON_TENANT_DB_DIR", "../tenant_db")
	return filepath.Join(tenantDBDir, fmt.Sprintf("%d.db", id))
}



// テナントDBに接続する
func connectToTenantDBInSQLite(id int64) (*sqlx.DB, error) {
	p := tenantDBPath(id)

	db, err := sqlx.Open(sqliteDriverName, fmt.Sprintf("file:%s?mode=rw", p))
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant DB: %w", err)
	}
	return db, nil
}

var (
	db *sqlx.DB
)


func MigrateFromSqlite() {
	num := 100
	for i := 0; i <= num; i++ {
		tenantDB, err := connectToTenantDBInSQLite(int64(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer tenantDB.Close()

		db, err = connectToTenantDB(int64(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer db.Close()

		// players
		var players []PlayerRow
		if err = tenantDB.Select(
		&players,
		"SELECT * FROM player WHERE tenant_id=?",
		int64(i),
		); err != nil {
			fmt.Errorf("error Select player: %w", err)
		}
		if players != nil {
			if _, err := db.NamedExec(
				"INSERT INTO player (id, tenant_id, display_name, is_disqualified, created_at, updated_at) VALUES (:id, :tenant_id, :display_name, :is_disqualified, :created_at, :updated_at)",
				players,
			); err != nil {
				fmt.Errorf(
					"error Insert player at tenantDB: id=%s, displayName=%s, isDisqualified=%t, createdAt=%d, updatedAt=%d, %w",
					players[0].ID, players[0].DisplayName, false, players[0].CreatedAt, players[0].UpdatedAt, err,
				)
			}
		}


		// competition
		var competitions []CompetitionRow
		if err = tenantDB.Select(
		&competitions,
		"SELECT * FROM competition WHERE tenant_id=?",
		int64(i),
		); err != nil {
			fmt.Errorf("error Select player: %w", err)
		}

		if competitions != nil {
			if _, err := db.NamedExec(
				"INSERT INTO competition (id, tenant_id, title, finished_at, created_at, updated_at) VALUES (:id, :tenant_id, :title, :finished_at, :created_at, :updated_at)",
				competitions,
			); err != nil {
				fmt.Errorf(
					"error Insert player at tenantDB: id=%s, createdAt=%d, updatedAt=%d, %w",
					competitions[0].ID, competitions[0].CreatedAt, competitions[0].UpdatedAt, err,
				)
			}
		}

		// player_score
		var playerScores []PlayerScoreRow
		if err = tenantDB.Select(
		&playerScores,
		"SELECT * FROM player_score WHERE tenant_id=?",
		int64(i),
		); err != nil {
			fmt.Errorf("error Select player: %w", err)
		}

		if playerScores != nil {
			if _, err := db.NamedExec(
				"INSERT INTO player_score (id, tenant_id, player_id, competition_id, score, row_num ,created_at, updated_at) VALUES (:id, :tenant_id, :player_id, :completition_id, :score, :row_num, :created_at, :updated_at)",
				playerScores,
			); err != nil {
				fmt.Errorf("error Insert player_score")
			}
		}
		fmt.Println(i)
	}
}
