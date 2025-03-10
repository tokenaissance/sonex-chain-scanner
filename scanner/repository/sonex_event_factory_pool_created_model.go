package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventFactoryPoolCreatedModel = (*customSonexEventFactoryPoolCreatedModel)(nil)

type (
	// SonexEventFactoryPoolCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventFactoryPoolCreatedModel.
	SonexEventFactoryPoolCreatedModel interface {
		sonexEventFactoryPoolCreatedModel
		withSession(session sqlx.Session) SonexEventFactoryPoolCreatedModel
	}

	customSonexEventFactoryPoolCreatedModel struct {
		*defaultSonexEventFactoryPoolCreatedModel
	}
)

// NewSonexEventFactoryPoolCreatedModel returns a model for the database table.
func NewSonexEventFactoryPoolCreatedModel(conn sqlx.SqlConn) SonexEventFactoryPoolCreatedModel {
	return &customSonexEventFactoryPoolCreatedModel{
		defaultSonexEventFactoryPoolCreatedModel: newSonexEventFactoryPoolCreatedModel(conn),
	}
}

func (m *customSonexEventFactoryPoolCreatedModel) withSession(session sqlx.Session) SonexEventFactoryPoolCreatedModel {
	return NewSonexEventFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(session))
}
