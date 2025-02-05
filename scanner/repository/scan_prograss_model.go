package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScanPrograssModel = (*customScanPrograssModel)(nil)

type (
	// ScanPrograssModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScanPrograssModel.
	ScanPrograssModel interface {
		scanPrograssModel
		withSession(session sqlx.Session) ScanPrograssModel
		New(ctx context.Context, data *ScanPrograss) error
		UpdateHeight(ctx context.Context, newData *ScanPrograss) error
	}

	customScanPrograssModel struct {
		*defaultScanPrograssModel
	}
)

// NewScanPrograssModel returns a model for the database table.
func NewScanPrograssModel(conn sqlx.SqlConn) ScanPrograssModel {
	return &customScanPrograssModel{
		defaultScanPrograssModel: newScanPrograssModel(conn),
	}
}

func (m *customScanPrograssModel) withSession(session sqlx.Session) ScanPrograssModel {
	return NewScanPrograssModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customScanPrograssModel) New(ctx context.Context, data *ScanPrograss) error {
	query := fmt.Sprintf("INSERT INTO %s (created_at, updated_at, chain_id, height) VALUES ($1, $2, $3, $4)", m.table)
	_, err := m.conn.ExecCtx(ctx, query, NowUTC(), NowUTC(), data.ChainId, data.Height)
	return err
}

func (m *customScanPrograssModel) UpdateHeight(ctx context.Context, newData *ScanPrograss) error {
	query := fmt.Sprintf("UPDATE %s SET updated_at=$1, height=$2 WHERE id = $3", m.table)
	_, err := m.conn.ExecCtx(ctx, query, NowUTC(), newData.Height, newData.Id)
	return err
}
