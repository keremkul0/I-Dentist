package main

import (
	"dental-clinic-system/api/appointment"
	"dental-clinic-system/api/clinic"
	"dental-clinic-system/api/forgotPassword"
	"dental-clinic-system/api/login"
	"dental-clinic-system/api/logout"
	"dental-clinic-system/api/patient"
	"dental-clinic-system/api/procedure"
	"dental-clinic-system/api/resetPassword"
	"dental-clinic-system/api/role"
	"dental-clinic-system/api/sendEmail"
	"dental-clinic-system/api/signUpClinic"
	"dental-clinic-system/api/singUpUser"
	"dental-clinic-system/api/user"
	"dental-clinic-system/api/verifyEmail"
	"dental-clinic-system/application/appointmentService"
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/application/emailService"
	"dental-clinic-system/application/jwtService"
	"dental-clinic-system/application/loginService"
	"dental-clinic-system/application/passwordResetService"
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/application/roleService"
	"dental-clinic-system/application/signUpClinicService"
	"dental-clinic-system/application/singUpUserService"
	"dental-clinic-system/application/tokenService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/background-jobs"
	config2 "dental-clinic-system/infrastructure/config"
	"dental-clinic-system/infrastructure/mailDialer"
	"dental-clinic-system/infrastructure/postgres"
	redis2 "dental-clinic-system/infrastructure/redis"
	"dental-clinic-system/infrastructure/repository/appointmentRepository"
	"dental-clinic-system/infrastructure/repository/clinicRepository"
	"dental-clinic-system/infrastructure/repository/loginRepository"
	"dental-clinic-system/infrastructure/repository/passwordResetTokenRepository"
	"dental-clinic-system/infrastructure/repository/patientRepository"
	"dental-clinic-system/infrastructure/repository/procedureRepository"
	"dental-clinic-system/infrastructure/repository/redisRepository"
	"dental-clinic-system/infrastructure/repository/roleRepository"
	"dental-clinic-system/infrastructure/repository/tokenRepository"
	"dental-clinic-system/infrastructure/repository/userRepository"
	"dental-clinic-system/middleware/authMiddleware"
	"dental-clinic-system/middleware/contextTimeoutMiddleware"
	"dental-clinic-system/vault"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {

	configModel := config2.SetConfig("resources")

	// Initialize base logger with request ID support
	logger.InitBaseLogger(configModel.Log)

	// Get base logger for startup logs
	baseLogger := logger.GetBaseLogger()
	clientVault, err := vault.ConnectVault(configModel.Vault)
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("Error connecting to vault")
		panic("Error connecting to vault")
	}

	err = config2.ReadConfigFromVault(clientVault, configModel, baseLogger)
	if err != nil {
		baseLogger.Fatal().Err(err).Msg("Error reading config from vault")
		panic("Error reading config from vault")
	}

	db := postgres.ConnectDatabase(configModel.Database)
	Rdb := redis2.ConnectRedis(configModel.Redis)
	mailDialerInstance := mailDialer.SetupMailDialer(configModel.Email)

	postgres.MigrateDatabase(db)

	//helpers.SetJWTKey(configModel.JWT.SecretKey)

	//Repositories
	newClinicRepository := clinicRepository.NewRepository(db)
	newAppointmentRepository := appointmentRepository.NewRepository(db)
	newPatientRepository := patientRepository.NewRepository(db)
	newProcedureRepository := procedureRepository.NewRepository(db)
	newRoleRepository := roleRepository.NewRepository(db)
	newUserRepository := userRepository.NewRepository(db)
	newLoginRepository := loginRepository.NewRepository(db)
	newTokenRepository := tokenRepository.NewRepository(db)
	newPasswordResetTokenRepository := passwordResetTokenRepository.NewRepository(db)

	//Redis Repository
	newRedisRepository := redisRepository.NewRepository(Rdb)

	//Services
	newClinicService := clinicService.NewClinicService(newClinicRepository)
	newAppointmentService := appointmentService.NewAppointmentService(newAppointmentRepository)
	newPatientService := patientService.NewPatientService(newPatientRepository)
	newProcedureService := procedureService.NewProcedureService(newProcedureRepository)
	newRoleService := roleService.NewRoleService(newRoleRepository)
	newUserService := userService.NewUserService(newUserRepository, newRoleService)
	newLoginService := loginService.NewLoginService(newLoginRepository)
	newSignUpClinicService := signUpClinicService.NewSignUpClinicService(newClinicRepository, newUserRepository, newRedisRepository)
	newTokenService := tokenService.NewTokenService(newTokenRepository)
	newSignUpUserService := signUpUserService.NewSignUpUserService(newUserRepository, newRedisRepository, newUserService)
	newEmailService := emailService.NewEmailService(newUserRepository, newTokenRepository, mailDialerInstance)
	newJwtService := jwtService.NewJwtService(configModel.JWT.SecretKey)
	newPasswordResetService := passwordResetService.NewPasswordResetService(newEmailService, newPasswordResetTokenRepository, newUserRepository)

	//Handlers
	newClinicHandler := clinic.NewClinicHandlerController(newClinicService, newUserService, newRoleService, newJwtService)
	newAppointmentHandler := appointment.NewAppointmentHandler(newAppointmentService, newUserService, newPatientService, newJwtService)
	newPatientHandler := patient.NewPatientController(newPatientService, newUserService, newJwtService)
	newProcedureHandler := procedure.NewProcedureController(newProcedureService, newUserService, newRoleService, newJwtService)
	newRoleHandler := role.NewRoleController(newRoleService)
	newUserHandler := user.NewUserController(newUserService, newRoleService, newJwtService)
	newLoginHandler := login.NewLoginController(newLoginService, newJwtService)
	newSignUpClinicHandler := signUpClinic.NewSignUpClinicController(newSignUpClinicService)
	newSignUpUserHandler := singUpUser.NewSignUpUserHandler(newSignUpUserService)
	newLogoutHandler := logout.NewLogoutController(newTokenService)
	newVerifyEmailHandler := verifyEmail.NewVerifyEmailController(newEmailService, newJwtService)
	newSendEmailHandler := sendEmail.NewSendEmailController(newEmailService, newJwtService)
	newForgotPasswordHandler := forgotPassword.NewForgotPasswordController(newPasswordResetService)
	newResetPasswordHandler := resetPassword.NewResetPasswordController(newPasswordResetService)

	//Create a new Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	//Middlewares
	newAuthMiddleware := authMiddleware.NewAuthMiddleware(newTokenService, newJwtService)

	// Public routes (no authentication required)
	login.RegisterAuthRoutes(app, newLoginHandler)
	signUpClinic.RegisterSignupClinicRoutes(app, newSignUpClinicHandler)
	singUpUser.RegisterSignupUserRoutes(app, newSignUpUserHandler)
	verifyEmail.RegisterVerifyEmailRoutes(app, newVerifyEmailHandler)
	forgotPassword.RegisterForgotPasswordRoutes(app, newForgotPasswordHandler)
	resetPassword.RegisterResetPasswordRoutes(app, newResetPasswordHandler)

	// Create API group with authentication middleware
	api := app.Group("/api", newAuthMiddleware.Authenticate())

	// Register Secured Routes
	clinic.RegisterClinicRoutes(api, newClinicHandler)
	appointment.RegisterAppointmentRoutes(api, newAppointmentHandler)
	patient.RegisterPatientsRoutes(api, newPatientHandler)
	procedure.RegisterProcedureRoutes(api, newProcedureHandler)
	role.RegisterRoleRoutes(api, newRoleHandler)
	user.RegisterUserRoutes(api, newUserHandler)
	logout.RegisterLogoutRoutes(api, newLogoutHandler)
	sendEmail.RegisterSendEmailRoutes(api, newSendEmailHandler)

	//background services
	background_jobs.StartCleanExpiredJwtTokens(newTokenService)
	background_jobs.StartCleanExpiredPasswordResetTokens(newPasswordResetTokenRepository)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		baseLogger.Info().Msgf("Server started on port %d", configModel.Server.Port)
		if err := app.Listen(fmt.Sprintf(":%d", configModel.Server.Port)); err != nil {
			baseLogger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-quit
	baseLogger.Info().Msg("Closing signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Info().Msgf("The server could not be shut down: %v", err)
	}
	gracefulShutdown(ctx, server, db, Rdb, clientVault)
	baseLogger.Info().Msg("Successful shutdown of the server.")
}

func gracefulShutdown(app *fiber.App, db *gorm.DB, redis *redis.Client, vaultClient *api.Client) {
	baseLogger := logger.GetBaseLogger()

	baseLogger.Info().Msg("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		baseLogger.Error().Err(err).Msg("Failed to gracefully shutdown server")
	} else {
		baseLogger.Info().Msg("Server stopped gracefully.")
	}

	// Close database connection
	if db != nil {
		baseLogger.Info().Msg("Closing database connection...")
		sqlDB, err := db.DB()
		if err != nil {
			baseLogger.Error().Err(err).Msg("Failed to get sql.DB from gorm.DB")
		} else {
			if err := sqlDB.Close(); err != nil {
				baseLogger.Error().Err(err).Msg("Failed to close database connection")
			} else {
				baseLogger.Info().Msg("Database connection closed.")
			}
		}
	}

	// Close Redis connection
	if redis != nil {
		baseLogger.Info().Msg("Closing Redis connection...")
		if err := redis.Close(); err != nil {
			baseLogger.Error().Err(err).Msg("Failed to close Redis connection")
		} else {
			baseLogger.Info().Msg("Redis connection closed.")
		}
	}

	// Clear Vault client token
	if vaultClient != nil {
		baseLogger.Info().Msg("Clearing Vault client token...")
		vaultClient.ClearToken()
		baseLogger.Info().Msg("Vault client token cleared.")
	}

	baseLogger.Info().Msg("Server shutdown complete.")
}
