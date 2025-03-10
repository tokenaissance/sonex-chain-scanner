package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPositionManagerIncreaseLiquidityModel = (*customSonexEventPositionManagerIncreaseLiquidityModel)(nil)

type (
	// SonexEventPositionManagerIncreaseLiquidityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPositionManagerIncreaseLiquidityModel.
	SonexEventPositionManagerIncreaseLiquidityModel interface {
		sonexEventPositionManagerIncreaseLiquidityModel
		withSession(session sqlx.Session) SonexEventPositionManagerIncreaseLiquidityModel
	}

	customSonexEventPositionManagerIncreaseLiquidityModel struct {
		*defaultSonexEventPositionManagerIncreaseLiquidityModel
	}
)

// NewSonexEventPositionManagerIncreaseLiquidityModel returns a model for the database table.
func NewSonexEventPositionManagerIncreaseLiquidityModel(conn sqlx.SqlConn) SonexEventPositionManagerIncreaseLiquidityModel {
	return &customSonexEventPositionManagerIncreaseLiquidityModel{
		defaultSonexEventPositionManagerIncreaseLiquidityModel: newSonexEventPositionManagerIncreaseLiquidityModel(conn),
	}
}

func (m *customSonexEventPositionManagerIncreaseLiquidityModel) withSession(session sqlx.Session) SonexEventPositionManagerIncreaseLiquidityModel {
	return NewSonexEventPositionManagerIncreaseLiquidityModel(sqlx.NewSqlConnFromSession(session))
}
