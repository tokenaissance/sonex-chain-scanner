package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPoolBurnModel = (*customSonexEventPoolBurnModel)(nil)

type (
	// SonexEventPoolBurnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPoolBurnModel.
	SonexEventPoolBurnModel interface {
		sonexEventPoolBurnModel
		withSession(session sqlx.Session) SonexEventPoolBurnModel
	}

	customSonexEventPoolBurnModel struct {
		*defaultSonexEventPoolBurnModel
	}
)

// NewSonexEventPoolBurnModel returns a model for the database table.
func NewSonexEventPoolBurnModel(conn sqlx.SqlConn) SonexEventPoolBurnModel {
	return &customSonexEventPoolBurnModel{
		defaultSonexEventPoolBurnModel: newSonexEventPoolBurnModel(conn),
	}
}

func (m *customSonexEventPoolBurnModel) withSession(session sqlx.Session) SonexEventPoolBurnModel {
	return NewSonexEventPoolBurnModel(sqlx.NewSqlConnFromSession(session))
}
