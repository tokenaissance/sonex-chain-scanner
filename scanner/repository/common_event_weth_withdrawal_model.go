package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommonEventWethWithdrawalModel = (*customCommonEventWethWithdrawalModel)(nil)

type (
	// CommonEventWethWithdrawalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommonEventWethWithdrawalModel.
	CommonEventWethWithdrawalModel interface {
		commonEventWethWithdrawalModel
		withSession(session sqlx.Session) CommonEventWethWithdrawalModel
	}

	customCommonEventWethWithdrawalModel struct {
		*defaultCommonEventWethWithdrawalModel
	}
)

// NewCommonEventWethWithdrawalModel returns a model for the database table.
func NewCommonEventWethWithdrawalModel(conn sqlx.SqlConn) CommonEventWethWithdrawalModel {
	return &customCommonEventWethWithdrawalModel{
		defaultCommonEventWethWithdrawalModel: newCommonEventWethWithdrawalModel(conn),
	}
}

func (m *customCommonEventWethWithdrawalModel) withSession(session sqlx.Session) CommonEventWethWithdrawalModel {
	return NewCommonEventWethWithdrawalModel(sqlx.NewSqlConnFromSession(session))
}
