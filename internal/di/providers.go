package di

import (
	"fmt"
	"identify/internal/config"
	"identify/internal/database"
	"identify/internal/handler"
	"identify/internal/logger"
	"identify/internal/provider"
	"identify/internal/provider/http"
	"identify/internal/repository"
	"identify/internal/security"
	"identify/internal/service"
)

// BuildContainer wires all application dependencies and constructs the Container.
func BuildContainer() (*Container, error) {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Init logger
	log, err := logger.NewZapLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 3. Connect database
	db, err := database.ConnectPostgresDb(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 4. Initialize securities
	passwordEncoder := security.NewBcryptPasswordEncoder()
	jwtProvider := security.NewJwtProviderImp(
		cfg.Jwt.AccessSecretKey,
		cfg.Jwt.RefreshSecretKey,
		cfg.Jwt.AccessTokenTTL,
		cfg.Jwt.RefreshTokenTTL,
	)

	// 5. Build repositories
	userRepo := repository.NewUserRepository(db)

	// 5.5 Build transaction manager
	txManager := repository.NewTransactionManager(db)

	// 6. Build services
	userService := service.NewUserService(userRepo, passwordEncoder, txManager)
	authService := service.NewAuthService(userRepo, passwordEncoder, jwtProvider)

	// 7. Build handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// 8. Build HTTP provider
	httpProvider := http.NewHTTPProvider(cfg.App, userHandler, authHandler, log)

	// 9. Build Registry
	registry := provider.NewRegistry()
	registry.Register(httpProvider)

	return &Container{
		Config:   cfg,
		Logger:   log,
		DB:       db,
		Registry: registry,
	}, nil
}
