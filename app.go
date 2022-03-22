package api

import (
	"api/migrations"
	"api/modules/account"
	"api/modules/users"
	"api/providers/binding"
	"api/providers/config"
	"api/providers/db"
	"api/providers/db/migrator"
	"api/providers/jwt"
	"api/providers/log"
	"api/providers/vm"
	"api/routers"

	"github.com/go-flow/flow/v2"
)

// AppModule is application root module
type AppModule struct {
	AppConfig config.AppConfig
	Logger    log.Logger
	Store     db.Store
	Injector  flow.Injector
}

// Start -
func (app *AppModule) Start() error {

	app.Logger.Info("Start application migrations...")
	// execute migration
	fsm := migrator.NewFSMigrator(migrations.Data, "mysql", app.Store)
	if err := fsm.Up(); err != nil {
		app.Logger.Fatal(err)
		return err
	}
	app.Logger.Info("End application migrations.")

	app.Logger.Infof("Application is running on: %s", app.Options().Addr)

	return nil
}

func (app *AppModule) Options() flow.Options {
	opts := flow.NewOptions()
	opts.Name = "core-api"

	addr := app.AppConfig.Addr()
	if addr != "" {
		opts.Addr = app.AppConfig.Addr()
	}

	return opts
}

func (app *AppModule) ProvideImports() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(config.New),
		flow.NewProvider(log.New),
		flow.NewProvider(db.New),
		flow.NewProvider(binding.New),
		flow.NewProvider(vm.NewJson),
		flow.NewProvider(jwt.NewAuth),
	}
}

func (app *AppModule) ProvideExports() []flow.Provider {
	return []flow.Provider{}
}

func (app *AppModule) ProvideModules() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(account.NewModule),
		flow.NewProvider(users.NewModule),
	}
}

// ProvideRouters handle http routing
func (app *AppModule) ProvideRouters() []flow.Provider {
	return []flow.Provider{
		flow.NewProvider(routers.NewRouter),
	}
}
