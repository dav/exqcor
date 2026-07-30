package store

type Script struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Theme          string `json:"theme"`
	WritingSeconds int    `json:"writing_seconds"`
	StationMode    string `json:"station_mode"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type Actor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type Character struct {
	ID          int64  `json:"id"`
	ScriptID    int64  `json:"script_id"`
	ActorID     *int64 `json:"actor_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type Section struct {
	ID          int64  `json:"id"`
	ScriptID    int64  `json:"script_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ordering    int    `json:"ordering"`
	Status      string `json:"status"`
}

type CharacterSection struct {
	ID          int64 `json:"id"`
	CharacterID int64 `json:"character_id"`
	SectionID   int64 `json:"section_id"`
	OnStage     bool  `json:"on_stage"`
}

type Prop struct {
	ID          int64  `json:"id"`
	SectionID   int64  `json:"section_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubSection struct {
	ID          int64   `json:"id"`
	SectionID   int64   `json:"section_id"`
	Ordering    int     `json:"ordering"`
	WriterID    *int64  `json:"writer_id"`
	StartedAt   *string `json:"started_at"`
	EndsAt      *string `json:"ends_at"`
	CompletedAt *string `json:"completed_at"`
}

type Line struct {
	ID           int64  `json:"id"`
	SubSectionID int64  `json:"sub_section_id"`
	CharacterID  int64  `json:"character_id"`
	Text         string `json:"text"`
	Ordering     int    `json:"ordering"`
	CreatedAt    string `json:"created_at"`
}

type Writer struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	AudienceMemberID *int64 `json:"audience_member_id"`
}

type AudienceMember struct {
	ID          int64   `json:"id"`
	ScriptID    int64   `json:"script_id"`
	Number      int     `json:"number"`
	Name        string  `json:"name"`
	DeviceToken string  `json:"-"`
	Status      string  `json:"status"`
	CalledAt    *string `json:"called_at"`
	CreatedAt   string  `json:"created_at"`
}
