//+build wireinject

package app

import "github.com/google/wire"

func InitializeApp(cfgFile string) (*AccountApp, error) {
	wire.Build(NewAccountApp, ReadConfig)
	return &AccountApp{}, nil
}
