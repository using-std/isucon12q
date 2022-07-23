package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"

	// "github.com/mattn/go-sqlite3"
	// proxy "github.com/shogo82148/go-sql-proxy"
) 

var (
	// 正しいテナント名の正規表現
	tenantNameRegexp = regexp.MustCompile(`^[a-z][a-z0-9-]{0,61}[a-z0-9]$`)

	adminDB *sqlx.DB

	sqliteDriverName = "sqlite3"
)

const (
	tenantDBSchemaFilePath  = "../sql/tenant/10_schema.sql"
	tenantDBSchemaFilePath2 = "../sql/tenant/11_add_index.sql"
	initializeScript        = "../sql/init.sh"
	cookieName              = "isuports_session"

	RoleAdmin     = "admin"
	RoleOrganizer = "organizer"
	RolePlayer    = "player"
	RoleNone      = "none"
)

// 環境変数を取得する、なければデフォルト値を返す
func getEnv(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

// テナントDBのパスを返す
func tenantDBPath(id int64) string {
	tenantDBDir := getEnv("ISUCON_TENANT_DB_DIR", "../tenant_db")
	return filepath.Join(tenantDBDir, fmt.Sprintf("%d.db", id))
}

// テナントDBに接続する
func connectToTenantDB(id int64) (*sqlx.DB, error) {
	p := tenantDBPath(id)

	db, err := sqlx.Open(sqliteDriverName, fmt.Sprintf("file:%s?mode=rw", p))
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant DB: %w", err)
	}
	return db, nil
}

// 管理用DBに接続する
func connectAdminDB() (*sqlx.DB, error) {
	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = getEnv("ISUCON_DB_HOST", "127.0.0.1") + ":" + getEnv("ISUCON_DB_PORT", "3306")
	config.User = getEnv("ISUCON_DB_USER", "isucon")
	config.Passwd = getEnv("ISUCON_DB_PASSWORD", "isucon")
	config.DBName = getEnv("ISUCON_DB_NAME", "isuports")
	config.ParseTime = true
	config.InterpolateParams = true
	dsn := config.FormatDSN()
	return sqlx.Open("mysql", dsn)
}

type PlayerRow struct {
	TenantID       int64  `db:"tenant_id"`
	ID             string `db:"id"`
	DisplayName    string `db:"display_name"`
	IsDisqualified bool   `db:"is_disqualified"`
	CreatedAt      int64  `db:"created_at"`
	UpdatedAt      int64  `db:"updated_at"`
}

type CompetitionRow struct {
	TenantID   int64         `db:"tenant_id"`
	ID         string        `db:"id"`
	Title      string        `db:"title"`
	FinishedAt sql.NullInt64 `db:"finished_at"`
	CreatedAt  int64         `db:"created_at"`
	UpdatedAt  int64         `db:"updated_at"`
}

type PlayerScoreRow struct {
	TenantID      int64  `db:"tenant_id"`
	ID            string `db:"id"`
	PlayerID      string `db:"player_id"`
	CompetitionID string `db:"competition_id"`
	Score         int64  `db:"score"`
	RowNum        int64  `db:"row_num"`
	CreatedAt     int64  `db:"created_at"`
	UpdatedAt     int64  `db:"updated_at"`
}

func main() {
	num := 100
	for i := 0; i <= num; i++ {
		tenantDB, err := connectToTenantDB(int64(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer tenantDB.Close()

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
			if _, err := tenantDB.NamedExec(
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
			if _, err := tenantDB.NamedExec(
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
			if _, err := tenantDB.NamedExec(
				"INSERT INTO player_score (id, tenant_id, player_id, competition_id, score, row_num ,created_at, updated_at) VALUES (:id, :tenant_id, :player_id, :completition_id, :score, :row_num, :created_at, :updated_at)",
				playerScores,
			); err != nil {
				fmt.Errorf("error Insert player_score")
			}
		}
		fmt.Println(i)
	}
}
