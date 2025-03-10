package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPositionManagerCollectModel = (*customSonexEventPositionManagerCollectModel)(nil)

type (
	// SonexEventPositionManagerCollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPositionManagerCollectModel.
	SonexEventPositionManagerCollectModel interface {
		sonexEventPositionManagerCollectModel
		withSession(session sqlx.Session) SonexEventPositionManagerCollectModel
	}

	customSonexEventPositionManagerCollectModel struct {
		*defaultSonexEventPositionManagerCollectModel
	}
)

// NewSonexEventPositionManagerCollectModel returns a model for the database table.
func NewSonexEventPositionManagerCollectModel(conn sqlx.SqlConn) SonexEventPositionManagerCollectModel {
	return &customSonexEventPositionManagerCollectModel{
		defaultSonexEventPositionManagerCollectModel: newSonexEventPositionManagerCollectModel(conn),
	}
}

func (m *customSonexEventPositionManagerCollectModel) withSession(session sqlx.Session) SonexEventPositionManagerCollectModel {
	return NewSonexEventPositionManagerCollectModel(sqlx.NewSqlConnFromSession(session))
}
