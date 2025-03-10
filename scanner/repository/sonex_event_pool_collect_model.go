package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPoolCollectModel = (*customSonexEventPoolCollectModel)(nil)

type (
	// SonexEventPoolCollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPoolCollectModel.
	SonexEventPoolCollectModel interface {
		sonexEventPoolCollectModel
		withSession(session sqlx.Session) SonexEventPoolCollectModel
	}

	customSonexEventPoolCollectModel struct {
		*defaultSonexEventPoolCollectModel
	}
)

// NewSonexEventPoolCollectModel returns a model for the database table.
func NewSonexEventPoolCollectModel(conn sqlx.SqlConn) SonexEventPoolCollectModel {
	return &customSonexEventPoolCollectModel{
		defaultSonexEventPoolCollectModel: newSonexEventPoolCollectModel(conn),
	}
}

func (m *customSonexEventPoolCollectModel) withSession(session sqlx.Session) SonexEventPoolCollectModel {
	return NewSonexEventPoolCollectModel(sqlx.NewSqlConnFromSession(session))
}
