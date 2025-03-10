package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPoolMintModel = (*customSonexEventPoolMintModel)(nil)

type (
	// SonexEventPoolMintModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPoolMintModel.
	SonexEventPoolMintModel interface {
		sonexEventPoolMintModel
		withSession(session sqlx.Session) SonexEventPoolMintModel
	}

	customSonexEventPoolMintModel struct {
		*defaultSonexEventPoolMintModel
	}
)

// NewSonexEventPoolMintModel returns a model for the database table.
func NewSonexEventPoolMintModel(conn sqlx.SqlConn) SonexEventPoolMintModel {
	return &customSonexEventPoolMintModel{
		defaultSonexEventPoolMintModel: newSonexEventPoolMintModel(conn),
	}
}

func (m *customSonexEventPoolMintModel) withSession(session sqlx.Session) SonexEventPoolMintModel {
	return NewSonexEventPoolMintModel(sqlx.NewSqlConnFromSession(session))
}
