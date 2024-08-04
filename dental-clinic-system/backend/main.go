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
    routes.RegisterRoutes(router,patientHandler)

    // Clinic Routes
    router.HandleFunc("/clinics", clinicHandler.GetClinics).Methods("GET")
    router.HandleFunc("/clinics/{id}", clinicHandler.GetClinic).Methods("GET")
    router.HandleFunc("/clinics", clinicHandler.CreateClinic).Methods("POST")
    router.HandleFunc("/clinics/{id}", clinicHandler.UpdateClinic).Methods("PUT")
    router.HandleFunc("/clinics/{id}", clinicHandler.DeleteClinic).Methods("DELETE")

    // Doctor Routes
    router.HandleFunc("/doctors", doctorHandler.GetDoctors).Methods("GET")
    router.HandleFunc("/doctors/{id}", doctorHandler.GetDoctor).Methods("GET")
    router.HandleFunc("/doctors", doctorHandler.CreateDoctor).Methods("POST")
    router.HandleFunc("/doctors/{id}", doctorHandler.UpdateDoctor).Methods("PUT")
    router.HandleFunc("/doctors/{id}", doctorHandler.DeleteDoctor).Methods("DELETE")

    // Assistant Routes
    router.HandleFunc("/assistants", assistantHandler.GetAssistants).Methods("GET")
    router.HandleFunc("/assistants/{id}", assistantHandler.GetAssistant).Methods("GET")
    router.HandleFunc("/assistants", assistantHandler.CreateAssistant).Methods("POST")
    router.HandleFunc("/assistants/{id}", assistantHandler.UpdateAssistant).Methods("PUT")
    router.HandleFunc("/assistants/{id}", assistantHandler.DeleteAssistant).Methods("DELETE")

    // Secretary Routes
    router.HandleFunc("/secretaries", secretaryHandler.GetSecretaries).Methods("GET")
    router.HandleFunc("/secretaries/{id}", secretaryHandler.GetSecretary).Methods("GET")
    router.HandleFunc("/secretaries", secretaryHandler.CreateSecretary).Methods("POST")
    router.HandleFunc("/secretaries/{id}", secretaryHandler.UpdateSecretary).Methods("PUT")
    router.HandleFunc("/secretaries/{id}", secretaryHandler.DeleteSecretary).Methods("DELETE")

    // Role Routes
    router.HandleFunc("/roles", roleHandler.GetRoles).Methods("GET")
    router.HandleFunc("/roles", roleHandler.CreateRole).Methods("POST")

    // Procedure Routes
    router.HandleFunc("/procedures", procedureHandler.GetProcedures).Methods("GET")
    router.HandleFunc("/procedures/{id}", procedureHandler.GetProcedure).Methods("GET")
    router.HandleFunc("/procedures", procedureHandler.CreateProcedure).Methods("POST")
    router.HandleFunc("/procedures/{id}", procedureHandler.UpdateProcedure).Methods("PUT")
    router.HandleFunc("/procedures/{id}", procedureHandler.DeleteProcedure).Methods("DELETE")

    // User Routes
    router.HandleFunc("/register", userHandler.Register).Methods("POST")
    router.HandleFunc("/login", userHandler.Login).Methods("POST")

    // Appointment Routes
    router.HandleFunc("/appointments", appointmentHandler.GetAppointments).Methods("GET")
    router.HandleFunc("/appointments/{id}", appointmentHandler.GetAppointment).Methods("GET")
    router.HandleFunc("/appointments", appointmentHandler.CreateAppointment).Methods("POST")
    router.HandleFunc("/appointments/{id}", appointmentHandler.UpdateAppointment).Methods("PUT")
    router.HandleFunc("/appointments/{id}", appointmentHandler.DeleteAppointment).Methods("DELETE")

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
