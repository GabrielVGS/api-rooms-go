package dtos

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Capacity    int    `json:"capacity"`
}

type RoomResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Subject     string            `json:"subject"`
	Capacity    int               `json:"capacity"`
	CreatedBy   uint              `json:"created_by"`
	Members     []RoomMemberResponse `json:"members,omitempty"`
	Notes       []NoteResponse    `json:"notes,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type UpdateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Capacity    int    `json:"capacity"`
}

type JoinRoomRequest struct {
	RoomID uint `json:"room_id"`
}

type RoomMemberResponse struct {
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
}
