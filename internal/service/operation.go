package service

import (
	"context"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
	"github.com/zenorachi/balance-management/internal/repository"
)

type OperationService struct {
	repo repository.Operation
}

func NewOperation(repo repository.Operation) *OperationService {
	return &OperationService{repo: repo}
}

func (o *OperationService) GetReportForUser(ctx context.Context, accountId int) ([]entity.Operation, error) {
	operations, err := o.getAllByAccountID(ctx, accountId)
	if err != nil {
		return nil, err
	}

	if len(operations) == 0 {
		return nil, entity.ErrAccountDoesNotExist
	}

	description := "%s money for product#%d from account#%d on %s"
	for i := 0; i < len(operations); i++ {
		if operations[i].OperationType == entity.TypeRevenue {
			operations[i].Description =
				fmt.Sprintf(description, "withdrawn", operations[i].OrderID, operations[i].AccountID, operations[i].OrderDate)
		} else {
			operations[i].Description =
				fmt.Sprintf(description, "credited", operations[i].OrderID, operations[i].AccountID, operations[i].OrderDate)
		}
	}

	return operations, nil
}

func (o *OperationService) GetReportForAccounting(ctx context.Context) ([]entity.Operation, error) {
	return o.repo.GetAll(ctx)
}

func (o *OperationService) getAllByAccountID(ctx context.Context, accountId int) ([]entity.Operation, error) {
	operations, err := o.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var operationsByID []entity.Operation
	for _, operation := range operations {
		if operation.AccountID == accountId {
			operationsByID = append(operationsByID, operation)
		}
	}

	return operationsByID, nil
}
