package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPoolSwapModel = (*customSonexEventPoolSwapModel)(nil)

type (
	// SonexEventPoolSwapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPoolSwapModel.
	SonexEventPoolSwapModel interface {
		sonexEventPoolSwapModel
		withSession(session sqlx.Session) SonexEventPoolSwapModel
	}

	customSonexEventPoolSwapModel struct {
		*defaultSonexEventPoolSwapModel
	}
)

// NewSonexEventPoolSwapModel returns a model for the database table.
func NewSonexEventPoolSwapModel(conn sqlx.SqlConn) SonexEventPoolSwapModel {
	return &customSonexEventPoolSwapModel{
		defaultSonexEventPoolSwapModel: newSonexEventPoolSwapModel(conn),
	}
}

func (m *customSonexEventPoolSwapModel) withSession(session sqlx.Session) SonexEventPoolSwapModel {
	return NewSonexEventPoolSwapModel(sqlx.NewSqlConnFromSession(session))
}
