package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommonEventWethDepositModel = (*customCommonEventWethDepositModel)(nil)

type (
	// CommonEventWethDepositModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommonEventWethDepositModel.
	CommonEventWethDepositModel interface {
		commonEventWethDepositModel
		withSession(session sqlx.Session) CommonEventWethDepositModel
	}

	customCommonEventWethDepositModel struct {
		*defaultCommonEventWethDepositModel
	}
)

// NewCommonEventWethDepositModel returns a model for the database table.
func NewCommonEventWethDepositModel(conn sqlx.SqlConn) CommonEventWethDepositModel {
	return &customCommonEventWethDepositModel{
		defaultCommonEventWethDepositModel: newCommonEventWethDepositModel(conn),
	}
}

func (m *customCommonEventWethDepositModel) withSession(session sqlx.Session) CommonEventWethDepositModel {
	return NewCommonEventWethDepositModel(sqlx.NewSqlConnFromSession(session))
}
