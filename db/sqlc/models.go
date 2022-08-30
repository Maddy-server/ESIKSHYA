// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

type Book struct {
	ID           int32          `json:"id"`
	BookName     string         `json:"book_name"`
	Content      sql.NullString `json:"content"`
	Writer       string         `json:"writer"`
	Section      string         `json:"section"`
	Randomunique string         `json:"randomunique"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}

type BookCount struct {
	ID        int32        `json:"id"`
	BookID    int32        `json:"book_id"`
	Count     int32        `json:"count"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type BookHistory struct {
	ID        int32     `json:"id"`
	BookID    int32     `json:"book_id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type BookRating struct {
	ID        int32        `json:"id"`
	BookID    int32        `json:"book_id"`
	Rating    int32        `json:"rating"`
	UserID    int32        `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type BookSaved struct {
	ID        int32     `json:"id"`
	BookID    int32     `json:"book_id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Child struct {
	ID         int32        `json:"id"`
	ParentID   int32        `json:"parent_id"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	Isverified sql.NullBool `json:"isverified"`
	CreatedAt  sql.NullTime `json:"created_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

type ChildNotification struct {
	ID              int32         `json:"id"`
	UserID          int32         `json:"user_id"`
	SecondaryUserID sql.NullInt32 `json:"secondary_user_id"`
	Title           string        `json:"title"`
	Type            string        `json:"type"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `json:"created_at"`
}

type ChildToken struct {
	ID     int32  `json:"id"`
	UserID int32  `json:"user_id"`
	Token  string `json:"token"`
}

type ChildrenDetail struct {
	ID          int32          `json:"id"`
	ChildrenID  int32          `json:"children_id"`
	FullName    string         `json:"full_name"`
	Grade       int32          `json:"grade"`
	DateOfBirth string         `json:"date_of_birth"`
	Gender      string         `json:"gender"`
	School      string         `json:"school"`
	Country     sql.NullString `json:"country"`
	State       sql.NullString `json:"state"`
}

type Friend struct {
	ID         int32        `json:"id"`
	SenderID   int32        `json:"sender_id"`
	ReceiverID int32        `json:"receiver_id"`
	Status     string       `json:"status"`
	FriendsAt  sql.NullTime `json:"friends_at"`
}

type FriendsLobbyQuestion struct {
	ID             int32  `json:"id"`
	LobbyID        int32  `json:"lobby_id"`
	Questions      string `json:"questions"`
	OptionsA       string `json:"options_a"`
	OptionsB       string `json:"options_b"`
	OptionsC       string `json:"options_c"`
	OptionsD       string `json:"options_d"`
	CorrectOptions string `json:"correct_options"`
}

type GameFriendLobby struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	OpID      int32     `json:"op_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GameNotification struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	OponentID   int32     `json:"oponent_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Status      string    `json:"status"`
	Grade       int32     `json:"grade"`
	CreatedAt   time.Time `json:"created_at"`
}

type GamePoint struct {
	ID           int32        `json:"id"`
	Player1ID    int32        `json:"player1_id"`
	Player2ID    int32        `json:"player2_id"`
	Player1Point int32        `json:"player1_point"`
	Player2Point int32        `json:"player2_point"`
	Indicator    int32        `json:"indicator"`
	PlayedTime   time.Time    `json:"played_time"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
}

type GameQuestion struct {
	ID              int32  `json:"id"`
	Class           int32  `json:"class"`
	Subject         string `json:"subject"`
	Questions       string `json:"questions"`
	OptionsA        string `json:"options_a"`
	OptionsB        string `json:"options_b"`
	OptionsC        string `json:"options_c"`
	OptionsD        string `json:"options_d"`
	CorrectOptions  string `json:"correct_options"`
	DifficultyLevel int32  `json:"difficulty_level"`
}

type GameQueue struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Status    string    `json:"status"`
	Subject   string    `json:"subject"`
	Grade     int32     `json:"grade"`
	LobbyID   int32     `json:"lobby_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GameRandomLobby struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Class     int32     `json:"class"`
	OpID      int32     `json:"op_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GeneralVideo struct {
	ID        int32          `json:"id"`
	Topic     sql.NullString `json:"topic"`
	Url       string         `json:"url"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

type Parent struct {
	ID         int32          `json:"id"`
	Phone      string         `json:"phone"`
	Password   sql.NullString `json:"password"`
	Otp        sql.NullString `json:"otp"`
	Isverified sql.NullBool   `json:"isverified"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	DeletedAt  sql.NullTime   `json:"deleted_at"`
}

type ParentsDetail struct {
	ID        int32  `json:"id"`
	RandomKey string `json:"random_key"`
	ParentID  int32  `json:"parent_id"`
	FullName  string `json:"full_name"`
	Address   string `json:"address"`
}

type ParentsNotification struct {
	ID              int32         `json:"id"`
	UserID          int32         `json:"user_id"`
	SecondaryUserID sql.NullInt32 `json:"secondary_user_id"`
	Title           string        `json:"title"`
	Type            string        `json:"type"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `json:"created_at"`
}

type ParentsToken struct {
	ID     int32  `json:"id"`
	UserID int32  `json:"user_id"`
	Token  string `json:"token"`
}

type Payment struct {
	ID               int32        `json:"id"`
	TransactionID    string       `json:"transaction_id"`
	TransactionToken string       `json:"transaction_token"`
	Method           string       `json:"method"`
	ParentID         int32        `json:"parent_id"`
	ChildID          int32        `json:"child_id"`
	Amount           int32        `json:"amount"`
	PayAt            sql.NullTime `json:"pay_at"`
	ExpireAt         sql.NullTime `json:"expire_at"`
}

type PaymentNumber struct {
	ID       int32        `json:"id"`
	Number   string       `json:"number"`
	Method   string       `json:"method"`
	Save     sql.NullBool `json:"save"`
	ParentID int32        `json:"parent_id"`
}

type RandomLobbyQuestion struct {
	ID             int32  `json:"id"`
	LobbyID        int32  `json:"lobby_id"`
	Questions      string `json:"questions"`
	OptionsA       string `json:"options_a"`
	OptionsB       string `json:"options_b"`
	OptionsC       string `json:"options_c"`
	OptionsD       string `json:"options_d"`
	CorrectOptions string `json:"correct_options"`
}

type Score struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	OwnPoints int32     `json:"own_points"`
	OpID      int32     `json:"op_id"`
	OpPoints  int32     `json:"op_points"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
}

type TimeTable struct {
	ID          int32          `json:"id"`
	ChildrenID  int32          `json:"children_id"`
	Class       int32          `json:"class"`
	Section     string         `json:"section"`
	Description string         `json:"description"`
	Day         string         `json:"day"`
	StartTime   sql.NullString `json:"start_time"`
	EndTime     sql.NullString `json:"end_time"`
}

type Video struct {
	ID        int32          `json:"id"`
	Class     int32          `json:"class"`
	Subject   string         `json:"subject"`
	Topic     sql.NullString `json:"topic"`
	Url       string         `json:"url"`
	CreatedAt sql.NullTime   `json:"created_at"`
	ImgUrl    sql.NullString `json:"img_url"`
	VideoID   sql.NullString `json:"video_id"`
}
