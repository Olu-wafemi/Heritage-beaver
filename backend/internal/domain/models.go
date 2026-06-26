package domain

import "time"

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	DisplayName    string    `json:"display_name"`
	PasswordHash   string    `json:"-"`
	PrimaryCulture string    `json:"primary_culture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FamilyMember struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DisplayName     string     `json:"display_name"`
	Gender          string     `json:"gender"`
	BirthDate       *time.Time `json:"birth_date"`
	DeathDate       *time.Time `json:"death_date"`
	BirthPlace      string     `json:"birth_place"`
	Biography       string     `json:"biography"`
	IsLiving        bool       `json:"is_living"`
	PrimaryLanguage string     `json:"primary_language"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Relationship struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	SourceMemberID   string    `json:"source_member_id"`
	TargetMemberID   string    `json:"target_member_id"`
	RelationshipType string    `json:"relationship_type"`
	CreatedAt        time.Time `json:"created_at"`
}

type Story struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	FamilyMemberID *string    `json:"family_member_id,omitempty"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	SourceType     string     `json:"source_type"`
	SourceLanguage string     `json:"source_language"`
	Summary        string     `json:"summary"`
	OccurredAt     *time.Time `json:"occurred_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type WisdomExtract struct {
	ID         string    `json:"id"`
	StoryID    string    `json:"story_id"`
	Excerpt    string    `json:"excerpt"`
	WisdomType string    `json:"wisdom_type"`
	Language   string    `json:"language"`
	Meaning    string    `json:"meaning"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}

type MythChapter struct {
	ID          string
	UserID      string
	Title       string
	Theme       string
	ChapterType string
	Narrative   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AncestorProfile struct {
	ID              string
	UserID          string
	FamilyMemberID  *string
	Name            string
	PersonaSummary  string
	VoiceStyle      string
	GuidanceStyle   string
	CulturalContext string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type GuidanceSession struct {
	ID                string
	UserID            string
	AncestorProfileID *string
	MythChapterID     *string
	Prompt            string
	Response          string
	SessionType       string
	CreatedAt         time.Time
}

type MediaAsset struct {
	ID         string
	UserID     string
	AssetType  string
	StorageKey string
	SourceKind string
	MimeType   string
	CreatedAt  time.Time
}

type PrivacySetting struct {
	ID                      string
	UserID                  string
	AllowFamilySharing      bool
	AllowVoiceCloning       bool
	AllowAncestorSimulation bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
