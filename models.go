package backend

import "time"

type AuthResponse struct {
	UserId  int    `json:"user_id"`
	Fname   string `json:"fname"`
	Lname   string `json:"lname"`
	Nname   string `json:"nname"`
	Avatar  string `json:"avatar"`
	About   string `json:"about"`
	Email   string `json:"email"`
	Dob     string `json:"dob"`
	Success bool   `json:"success"`
	Public  int    `json:"public"`
}

type loginPayload struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

type regPayload struct {
	Email  string `json:"email"`
	Pw     string `json:"pw"`
	Fname  string `json:"fname"`
	Lname  string `json:"lname"`
	Dob    string `json:"dob"`
	Avatar string `json:"avatar"`
	Nname  string `json:"nname"`
	About  string `json:"about"`
}

type UserStruct struct {
	Id       int    `json:"id"`
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Nname    string `json:"nname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Dob      string `json:"dob"`
	Avatar   string `json:"avatar"`
	About    string `json:"about"`
	Public   int    `json:"public"`
}

type UserPayload struct {
	Data []UserStruct `json:"data"`
}
type UserPayloadWithStatus struct {
	Data   UserStruct `json:"data"`
	Status int        `json:"status"`
}

type UserId struct {
	TargetId int `json:"targetid"`
	Public   int `json:"public"`
}

type UserFollowerStruct struct {
	Id        int    `json:"id"`
	Action    string `json:"action"`
	SourceId  int    `json:"sourceid"`
	TargetId  int    `json:"targetid"`
	Status    int    `json:"status"`
	ChatNoti  int    `json:"chat_noti"`
	LastMsgAt string `json:"last_msg_at"`
}

type UserMessageStruct struct {
	Label         string `json:"label"`
	Id            int    `json:"id"`
	TargetId      int    `json:"targetid"`
	SourceId      int    `json:"sourceid"`
	Message       string `json:"message"`
	CreatedAt     string `json:"createdat"`
	OnlineUserIds []int  `json:"onlineuserids"`
}

type UserMessagePayload struct {
	Data []UserMessageStruct `json:"data"`
}

type PostStruct struct {
	Id        int    `json:"id"`
	Author    int    `json:"author"` // author uid
	Message   string `json:"message"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdat"`
	Privacy   int    `json:"privacy"`
}

type PostResponse struct {
	Id        int    `json:"id"`
	Author    int    `json:"author"` // author uid
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	Avatar    string `json:"avatar"`
	Nname     string `json:"nname"`
	Message   string `json:"message"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdat"`
	Success   bool   `json:"success"`
	Privacy   int    `json:"privacy"`
}

type PostPayload struct {
	Data []PostResponse `json:"data"`
}

// type PostPayload struct {
// 	UserId  int    `json:"user_id"`
// 	Content string `json:"content"`
// 	Image   string `json:"image"`
// 	Privacy string `json:"privacy"`
// }

type PostMemberStruct struct {
	Id         int `json:"id"`
	UserId     int `json:"userid"`
	UserPostId int `json:"userpostid"`
}

type PostCommentStruct struct {
	Id        int       `json:"id"`
	PostId    int       `json:"postid"`
	UserId    int       `json:"userid"`
	CreatedAt time.Time `json:"createdat"`
	Message   string    `json:"message"`
	Image     string    `json:"image"`
}

type PostCommentResponse struct {
	Id        int    `json:"id"`
	PostId    int    `json:"postid"`
	UserId    int    `json:"userid"`
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	Avatar    string `json:"avatar"`
	Nname     string `json:"nname"`
	CreatedAt string `json:"createdat"`
	Message   string `json:"message"`
	Image     string `json:"image"`
	Success   bool   `json:"success"`
}

type PostCommentPayload struct {
	Data []PostCommentStruct `json:"data"`
}

type GroupStruct struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Creator     int    `json:"creator"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdat"`
}

type GroupResponse struct {
	Creator        int    `json:"creator"`
	CreatedAt      string `json:"createdat"`
	Success        bool   `json:"success"`
	CreatedGroupId int    `json:"createdid"`
}

type GroupPayload struct {
	Data []GroupStruct `json:"data"`
}

type GroupRequestStruct struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userid"`
	GroupId   int    `json:"groupid"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdat"`
}

type GroupRequestPayload struct {
	Data []GroupRequestStruct `json:"data"`
}

type GroupMemberStruct struct {
	Id       int `json:"id"`
	UserId   int `json:"userid"`
	GroupId  int `json:"groupid"`
	Status   int `json:"status"`
	ChatNoti int `json:"chatnoti"`
}

type GroupEventStruct struct {
	Id          int    `json:"id"`
	GroupId     int    `json:"groupid"`
	Author      int    `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdat"`
	Date        string `json:"date"`
}

type GroupEventPayload struct {
	Data []GroupEventStruct `json:"data"`
}

type GroupEventMemberStruct struct {
	Id      int `json:"id"`
	Status  int `json:"status"`
	UserId  int `json:"userid"`
	EventId int `json:"eventid"`
}

type GroupEventMemberPayload struct {
	Data []GroupEventMemberStruct `json:"data"`
}

type GroupPostStruct struct {
	Id        int    `json:"id"`
	Fname     string `json:"fname"`
	Lname     string `json:"lname"`
	Nickname  string `json:"nname"`
	Author    int    `json:"author"`
	GroupId   int    `json:"groupid"`
	Message   string `json:"message"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdat"`
}

type GroupPostPayload struct {
	Data []GroupPostStruct `json:"data"`
}

type GroupPostCommentStruct struct {
	Id          int    `json:"id"`
	Fname       string `json:"fname"`
	Lname       string `json:"lname"`
	Nickname    string `json:"nname"`
	GroupPostId int    `json:"postid"`
	Author      int    `json:"userid"`
	CreatedAt   string `json:"createdat"`
	Message     string `json:"message"`
	AuthorAvatar string `json:"avatar"`
}

type GroupPostCommentPayload struct {
	Data []GroupPostCommentStruct `json:"data"`
}

type GroupMessageStruct struct {
	Label     string `json:"label"`
	Id        int    `json:"id"`
	Message   string `json:"message"`
	SourceId  int    `json:"sourceid"`
	GroupId   int    `json:"groupid"`
	CreatedAt string `json:"createdat"`
}

type GroupMessagePayload struct {
	Data []GroupMessageStruct `json:"data"`
}

type SessionStruct struct {
	SessionToken string `json:"sessiontoken"`
	UserId       int    `json:"userid"`
}

type NotifStruct struct {
	Label     string `json:"label"`
	Id        int    `json:"id"`
	Type      string `json:"type"`
	SourceId  int    `json:"sourceid"`
	TargetId  int    `json:"targetid"`
	Accepted  bool   `json:"accepted"`
	CreatedAt string `json:"createdat"`
	GroupId   int    `json:"groupid"`
}

type NotiMessageStruct struct {
	Label      string `json:"label"`
	Id         int    `json:"id"`
	SourceId   int    `json:"sourceid"`
	TargetId   int    `json:"targetid"`
	Message    string `json:"message"`
	GroupId    int    `json:"groupid"`
	CreatedAt  string `json:"createdat"`
	Type       string `json:"type"`
	Accepted   bool   `json:"accepted"`
	GroupTitle string `json:"grouptitle"`
}

type PrivateChatItemStruct struct {
	Id        int    `json:"id"`
	SourceId  int    `json:"sourceid"`
	TargetId  int    `json:"targetid"`
	ChatNoti  int    `json:"chat_noti"`
	LastMsgAt string `json:"last_msg_at"`
}

type PrivateChatItemPayload struct {
	Data []PrivateChatItemStruct `json:"data"`
}

type GroupChatItemStruct struct {
	Id        int    `json:"id"`
	GroupId   int    `json:"groupid"`
	UserId    int    `json:"userid"`
	ChatNoti  int    `json:"chat_noti"`
	LastMsgAt string `json:"last_msg_at"`
}

type GroupChatItemPayload struct {
	Data []GroupChatItemStruct `json:"data"`
}
