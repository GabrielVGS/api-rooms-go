package dtos

type CreateNoteRequest struct {
	RoomID  uint   `json:"room_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	RoomID    uint   `json:"room_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	RoomName  string `json:"room_name,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
