package manager

import (
	"github.com/ProtocolONE/p1pay.api/database/dao"
	"github.com/ProtocolONE/p1pay.api/database/model"
	"go.uber.org/zap"
)

type CurrencyManager Manager

func InitCurrencyManager(database dao.Database, logger *zap.SugaredLogger) *CurrencyManager {
	return &CurrencyManager{Database: database, Logger: logger}
}

func (cm *CurrencyManager) FindByCodeInt(codeInt int) *model.Currency {
	c, err := cm.Database.Repository(TableCurrency).FindCurrencyById(codeInt)

	if err != nil {
		cm.Logger.Errorf("Query from table \"%s\" ended with error: %s", TableCurrency, err)
	}

	return c
}

func (cm *CurrencyManager) FindByName(name string) []*model.Currency {
	c, err := cm.Database.Repository(TableCurrency).FindCurrenciesByName(name)

	if err != nil {
		cm.Logger.Errorf("Query from table \"%s\" ended with error: %s", TableCurrency, err)
	}

	if c == nil {
		return []*model.Currency{}
	}

	return c
}

func (cm *CurrencyManager) FindAll(limit int, offset int) []*model.Currency {
	c, err := cm.Database.Repository(TableCurrency).FindAllCurrencies(limit, offset)

	if err != nil {
		cm.Logger.Errorf("Query from table \"%s\" ended with error: %s", TableCurrency, err)
	}

	if c == nil {
		return []*model.Currency{}
	}

	return c
}
