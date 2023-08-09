package main

import (
	"backend"
	db "backend/pkg/db/sqlite"
	"backend/pkg/websocket"
	"fmt"
	"net/http"
)

func main() {

	db.RunMigration()
	db.DbConnect()
	// db.RemoveMigration(m)
	// db.InsertMockUserData()
	// db.InsertMockPostData()

	hub := websocket.NewHub()
	go hub.Run()

	// exec.Command("xdg-open", "https://localhost:8080").Start()

	mux := http.NewServeMux()

	mux.Handle("/", backend.Homehandler())
	mux.Handle("/session", backend.SessionHandler())
	mux.Handle("/login", backend.Loginhandler())
	mux.Handle("/logout", backend.Logouthandler())
	mux.Handle("/reg", backend.Reghandler())
	mux.Handle("/user", backend.Userhandler())
	mux.Handle("/privacy", backend.PrivacyHandler())
	mux.Handle("/user-follower", backend.UserFollowerHandler())
	mux.Handle("/user-following", backend.UserFollowingHandler())
	mux.Handle("/user-follow-status", backend.UserFollowerStatusHandler())
	mux.Handle("/close-friend", backend.CloseFriendHandler())
	mux.Handle("/user-message", backend.UserMessageHandler())
	mux.Handle("/post", backend.Posthandler())
	mux.Handle("/post-comment", backend.PostCommentHandler())
	mux.Handle("/group", backend.Grouphandler())
	mux.Handle("/group-member", backend.GroupMemberHandler())
	mux.Handle("/group-request", backend.GroupRequestHandler())
	mux.Handle("/group-post", backend.GroupPostHandler())
	mux.Handle("/group-post-comment", backend.GroupPostCommentHandler())
	mux.Handle("/group-event", backend.GroupEventHandler())
	mux.Handle("/group-event-member", backend.GroupEventMemberHandler())
	mux.Handle("/group-message", backend.GroupMessageHandler())
	mux.Handle("/group-request-by-user", backend.GroupRequestByUserHandler())

	mux.Handle("/private-chat-item", backend.PrivateChatItemHandler())
	mux.Handle("/group-chat-item", backend.GroupChatItemHandler())
	mux.Handle("/group-chat-seen", backend.GroupChatSeenHandler())
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	fmt.Println("Starting server at port 8080")

	err1 := http.ListenAndServe(":8080", mux)
	if err1 != nil {
		fmt.Println(err1)
	}
}
