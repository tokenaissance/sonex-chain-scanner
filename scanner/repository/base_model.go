package repository

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func NewConn(connectString string) sqlx.SqlConn {
	return sqlx.NewSqlConn("pgx", connectString)
}

func IsTx(conn sqlx.SqlConn) bool {
	_, err := conn.RawDB()
	return err != nil
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
