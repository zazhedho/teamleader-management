package dto

import "time"

// TL Daily Activity DTOs
type TLActivityCreate struct {
	Date             time.Time `json:"date" binding:"required"`
	ActivityType     string    `json:"activity_type" binding:"required,oneof=canvassing pameran"`
	Kecamatan        *string   `json:"kecamatan" binding:"omitempty,min=2,max=100"`
	Desa             *string   `json:"desa" binding:"omitempty,min=2,max=100"`
	GpsLat           *float64  `json:"gps_lat" binding:"omitempty,min=-90,max=90"`
	GpsLng           *float64  `json:"gps_lng" binding:"omitempty,min=-180,max=180"`
	DurationHours    *float64  `json:"duration_hours" binding:"omitempty,min=0,max=24"`
	ProspectCount    int       `json:"prospect_count" binding:"min=0"`
	DealCount        int       `json:"deal_count" binding:"min=0"`
	MotorkuDownloads int       `json:"motorku_downloads" binding:"min=0"`
	Notes            *string   `json:"notes" binding:"omitempty,max=500"`
}

type TLActivityUpdate struct {
	Date             *time.Time `json:"date" binding:"omitempty"`
	ActivityType     *string    `json:"activity_type" binding:"omitempty,oneof=canvassing pameran"`
	Kecamatan        *string    `json:"kecamatan" binding:"omitempty,min=2,max=100"`
	Desa             *string    `json:"desa" binding:"omitempty,min=2,max=100"`
	GpsLat           *float64   `json:"gps_lat" binding:"omitempty,min=-90,max=90"`
	GpsLng           *float64   `json:"gps_lng" binding:"omitempty,min=-180,max=180"`
	DurationHours    *float64   `json:"duration_hours" binding:"omitempty,min=0,max=24"`
	ProspectCount    *int       `json:"prospect_count" binding:"omitempty,min=0"`
	DealCount        *int       `json:"deal_count" binding:"omitempty,min=0"`
	MotorkuDownloads *int       `json:"motorku_downloads" binding:"omitempty,min=0"`
	Notes            *string    `json:"notes" binding:"omitempty,max=500"`
}

// TL Attendance DTOs
type AttendanceRecord struct {
	SalesmanPersonId string `json:"salesman_person_id" binding:"required"`
	SalesmanName     string `json:"salesman_name" binding:"required,min=2,max=150"`
	Status           string `json:"status" binding:"required,oneof=hadir tidak_hadir"`
}

type TLAttendanceCreate struct {
	Date       time.Time          `json:"date" binding:"required"`
	Attendance []AttendanceRecord `json:"attendance" binding:"required,min=1,dive"`
}

type TLAttendanceUpdate struct {
	Date       *time.Time         `json:"date" binding:"omitempty"`
	Attendance []AttendanceRecord `json:"attendance" binding:"omitempty,min=1,dive"`
}

// TL Session (Coaching & Briefing) DTOs - MERGED
type TLSessionCreate struct {
	SessionType   string    `json:"session_type" binding:"required,oneof=coaching briefing"`
	Date          time.Time `json:"date" binding:"required"`
	Notes         *string   `json:"notes" binding:"omitempty,max=1000"`
	Attendees     []string  `json:"attendees" binding:"omitempty"`
	DurationHours *float64  `json:"duration_hours" binding:"omitempty,min=0,max=24"`
}

type TLSessionUpdate struct {
	SessionType   *string    `json:"session_type" binding:"omitempty,oneof=coaching briefing"`
	Date          *time.Time `json:"date" binding:"omitempty"`
	Notes         *string    `json:"notes" binding:"omitempty,max=1000"`
	Attendees     []string   `json:"attendees" binding:"omitempty"`
	DurationHours *float64   `json:"duration_hours" binding:"omitempty,min=0,max=24"`
}

// TL Training Participation DTOs
type TrainingParticipant struct {
	SalesmanPersonId string `json:"salesman_person_id" binding:"required"`
	SalesmanName     string `json:"salesman_name" binding:"required,min=2,max=150"`
	Status           string `json:"status" binding:"required,oneof=hadir tidak_hadir"`
}

type TLTrainingCreate struct {
	TrainingName string                `json:"training_name" binding:"required,min=3,max=200"`
	Date         time.Time             `json:"date" binding:"required"`
	Participants []TrainingParticipant `json:"participants" binding:"required,min=1,dive"`
}
