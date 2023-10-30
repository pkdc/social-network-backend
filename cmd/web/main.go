package main

import (
	"backend"
	db "backend/pkg/db/sqlite"
	"backend/pkg/websocket"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {

	// db.RunMigration()
	db.DbConnect()
	// db.RemoveMigration(m)
	// db.InsertMockUserData()
	// db.InsertMockPostData()

	hub := websocket.NewHub()
	go hub.Run()

	// exec.Command("xdg-open", "https://localhost:8080").Start()

	// Create a new CORS middleware
	corsHandler := cors.New(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"https://notfacebook.netlify.app"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()

	mux.Handle("/", corsHandler.Handler(backend.Homehandler()))
	mux.Handle("/session", corsHandler.Handler(backend.SessionHandler()))
	mux.Handle("/login", corsHandler.Handler(backend.Loginhandler()))
	mux.Handle("/logout", corsHandler.Handler(backend.Logouthandler()))
	mux.Handle("/reg", corsHandler.Handler(backend.Reghandler()))
	mux.Handle("/user", corsHandler.Handler(backend.Userhandler()))
	mux.Handle("/privacy", corsHandler.Handler(backend.PrivacyHandler()))
	mux.Handle("/user-follower", corsHandler.Handler(backend.UserFollowerHandler()))
	mux.Handle("/user-following", corsHandler.Handler(backend.UserFollowingHandler()))
	mux.Handle("/user-follow-status", corsHandler.Handler(backend.UserFollowerStatusHandler()))
	mux.Handle("/close-friend", corsHandler.Handler(backend.CloseFriendHandler()))
	mux.Handle("/user-message", corsHandler.Handler(backend.UserMessageHandler()))
	mux.Handle("/post", corsHandler.Handler(backend.Posthandler()))
	mux.Handle("/post-comment", corsHandler.Handler(backend.PostCommentHandler()))
	mux.Handle("/group", corsHandler.Handler(backend.Grouphandler()))
	mux.Handle("/group-member", corsHandler.Handler(backend.GroupMemberHandler()))
	mux.Handle("/group-request", corsHandler.Handler(backend.GroupRequestHandler()))
	mux.Handle("/group-post", corsHandler.Handler(backend.GroupPostHandler()))
	mux.Handle("/group-post-comment", corsHandler.Handler(backend.GroupPostCommentHandler()))
	mux.Handle("/group-event", corsHandler.Handler(backend.GroupEventHandler()))
	mux.Handle("/group-event-member", corsHandler.Handler(backend.GroupEventMemberHandler()))
	mux.Handle("/group-message", corsHandler.Handler(backend.GroupMessageHandler()))
	mux.Handle("/group-request-by-user", corsHandler.Handler(backend.GroupRequestByUserHandler()))

	mux.Handle("/private-chat-item", corsHandler.Handler(backend.PrivateChatItemHandler()))
	mux.Handle("/group-chat-item", corsHandler.Handler(backend.GroupChatItemHandler()))
	mux.Handle("/group-chat-seen", corsHandler.Handler(backend.GroupChatSeenHandler()))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// fmt.Println("Starting server at port 8080")

	// err1 := http.ListenAndServe(":8080", mux)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }

	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, mux)
    if err != nil {
        log.Fatal(err)
    }
}
