package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexEventPositionManagerDecreaseLiquidityModel = (*customSonexEventPositionManagerDecreaseLiquidityModel)(nil)

type (
	// SonexEventPositionManagerDecreaseLiquidityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexEventPositionManagerDecreaseLiquidityModel.
	SonexEventPositionManagerDecreaseLiquidityModel interface {
		sonexEventPositionManagerDecreaseLiquidityModel
		withSession(session sqlx.Session) SonexEventPositionManagerDecreaseLiquidityModel
	}

	customSonexEventPositionManagerDecreaseLiquidityModel struct {
		*defaultSonexEventPositionManagerDecreaseLiquidityModel
	}
)

// NewSonexEventPositionManagerDecreaseLiquidityModel returns a model for the database table.
func NewSonexEventPositionManagerDecreaseLiquidityModel(conn sqlx.SqlConn) SonexEventPositionManagerDecreaseLiquidityModel {
	return &customSonexEventPositionManagerDecreaseLiquidityModel{
		defaultSonexEventPositionManagerDecreaseLiquidityModel: newSonexEventPositionManagerDecreaseLiquidityModel(conn),
	}
}

func (m *customSonexEventPositionManagerDecreaseLiquidityModel) withSession(session sqlx.Session) SonexEventPositionManagerDecreaseLiquidityModel {
	return NewSonexEventPositionManagerDecreaseLiquidityModel(sqlx.NewSqlConnFromSession(session))
}
