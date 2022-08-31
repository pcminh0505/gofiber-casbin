package database

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func Casbin() *casbin.Enforcer {
	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(adminDB)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	e, err := casbin.NewEnforcer("config/restful_rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	// Add policy - One-time run
	if hasPolicy := e.HasPolicy("admin", "/api/admin/*", "(GET)|(POST)|(PUT)|(DELETE)"); !hasPolicy {
		e.AddPolicy("admin", "/api/admin/*", "(GET)|(POST)|(PUT)|(DELETE)")
	}
	if hasPolicy := e.HasPolicy("user", "/api/users/:id/*", "(GET)|(PUT)"); !hasPolicy {
		e.AddPolicy("user", "/api/users/:id/*", "(GET)|(PUT)")
	}

	e.LoadPolicy()
	return e
}
