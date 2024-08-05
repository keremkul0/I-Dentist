package main

import (
	"dental-clinic-system/handlers"
	"dental-clinic-system/models"
	"dental-clinic-system/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func init() {
    dsn := "host=localhost user=clinicuser password=clinicpassword dbname=clinicdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Automigrate models
    db.AutoMigrate(
        &models.Patient{}, 
        &models.Clinic{}, 
        &models.Doctor{}, 
        &models.Assistant{}, 
        &models.Secretary{}, 
        &models.Appointment{}, 
        &models.Role{}, 
        &models.User{}, 
        &models.Procedure{},
    )
}

func main() {
    router := mux.NewRouter()

    patientHandler := &handlers.PatientHandler{DB: db}
    clinicHandler := &handlers.ClinicHandler{DB: db}
    doctorHandler := &handlers.DoctorHandler{DB: db}
    assistantHandler := &handlers.AssistantHandler{DB: db}
    secretaryHandler := &handlers.SecretaryHandler{DB: db}
    appointmentHandler := &handlers.AppointmentHandler{DB: db}
    roleHandler := &handlers.RoleHandler{DB: db}
    procedureHandler := &handlers.ProcedureHandler{DB: db}
    userHandler := &handlers.UserHandler{DB: db}


    // Patient Routes
    routes.RegisterpatientsRoutes(router,patientHandler)

    // Clinic Routes
    routes.RegisterClinicRoutes(router,clinicHandler)

    // Doctor Routes
    routes.RegisterDoctorRoutes(router,doctorHandler)

    // Assistant Routes
    routes.RegisterAssistantRoutes(router,assistantHandler)

    // Secretary Routes
    routes.RegisterSecretaryRoutes(router,secretaryHandler)

    // Role Routes
    routes.RegisterRoleRoutes(router,roleHandler)

    // Procedure Routes
    routes.RegisterProcedureRoutes(router,procedureHandler)

    // User Routes
    routes.RegisterUserRoutes(router,userHandler)

    // Appointment Routes
    routes.RegisterAppointmentRoutes(router,appointmentHandler)
    

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
