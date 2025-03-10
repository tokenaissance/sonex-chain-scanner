package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPositionManagerTransferModel = (*customSonexEventPositionManagerTransferModel)(nil)

type (
	// SonexEventPositionManagerTransferModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPositionManagerTransferModel.
	SonexEventPositionManagerTransferModel interface {
		sonexEventPositionManagerTransferModel
		withSession(session sqlx.Session) SonexEventPositionManagerTransferModel
	}

	customSonexEventPositionManagerTransferModel struct {
		*defaultSonexEventPositionManagerTransferModel
	}
)

// NewSonexEventPositionManagerTransferModel returns a model for the database table.
func NewSonexEventPositionManagerTransferModel(conn sqlx.SqlConn) SonexEventPositionManagerTransferModel {
	return &customSonexEventPositionManagerTransferModel{
		defaultSonexEventPositionManagerTransferModel: newSonexEventPositionManagerTransferModel(conn),
	}
}

func (m *customSonexEventPositionManagerTransferModel) withSession(session sqlx.Session) SonexEventPositionManagerTransferModel {
	return NewSonexEventPositionManagerTransferModel(sqlx.NewSqlConnFromSession(session))
}
