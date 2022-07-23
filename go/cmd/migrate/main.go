package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jmoiron/sqlx"
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

func main() {
	for i := 0; i <= 100; i++ {
		db, err := connectToTenantDB(int64(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer db.Close()
		fmt.Println(db.DriverName())
	}
}
