package route

import (
	"itish41/doctor_ai_assistant/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Authentication routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", controller.SignUp)   //done
		authGroup.POST("/login", controller.Login)     //done
		authGroup.POST("/me", controller.GetProfile)   //done
		authGroup.PUT("/me", controller.UpdateProfile) //done
	}

	// Transcription routes
	transcriptionGroup := r.Group("/transcription")
	{
		transcriptionGroup.POST("/", controller.GetTranscriptions)                            //done
		transcriptionGroup.POST("/:id/download", controller.DownloadTranscription)            //done
		transcriptionGroup.POST("/:id/getTranscriptionByID", controller.GetTranscriptionByID) //done
		transcriptionGroup.PUT("/:id", controller.UpdateTranscription)                        //done
		transcriptionGroup.POST("/:id/delete", controller.DeleteTranscription)                // done
	}

	// Patient routes
	patientsGroup := r.Group("/patients")
	{
		patientsGroup.POST("/", controller.CreatePatient)                  //done
		patientsGroup.POST("/patientsList", controller.GetPatients)        //done
		patientsGroup.POST("/:id/getPatient", controller.GetPatientByID)   // done
		patientsGroup.PUT("/:id/update", controller.UpdatePatient)         // done
		patientsGroup.POST("/:id", controller.DeletePatient)               // done
		patientsGroup.POST("/transcript", controller.GetPatientTranscript) // Get patient transcript by name
	}

	// Dashboard & Statistics routes
	dashboardGroup := r.Group("/dashboard")
	{
		dashboardGroup.POST("/transcripts", controller.GetDashboardTranscripts)
		dashboardGroup.POST("/patients", controller.GetDashboardPatients)
		dashboardGroup.POST("/statistics/daily", controller.GetDailyStatistics)
		dashboardGroup.POST("/statistics/monthly", controller.GetMonthlyStatistics)
		dashboardGroup.POST("/statistics/busiest-days", controller.GetBusiestDays)
	}

	// Optionally, for audio upload (future feature)
	// r.POST("/upload_audio", controller.UploadAudio)
}
