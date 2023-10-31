package backend

import (
	"backend/pkg/db/crud"
	db "backend/pkg/db/sqlite"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func UrlPathMatcher(w http.ResponseWriter, r *http.Request, p string) error {
	if r.URL.Path != p {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return errors.New("404")
	}

	return nil
}

func WriteHttpHeader(jsonResp []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// func // EnableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
// }

func Homehandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
	}
}

func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/session"); err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Declares the payload struct
			var Resp SessionStruct

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET SESSION FOR USER ###

			session, err := r.Cookie("SessionToken")

			sessionTable, err := query.GetUserId(context.Background(), session.Value)

			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			Resp.UserId = int(sessionTable.UserID)
			Resp.SessionToken = sessionTable.SessionToken

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser

			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

	}
}

func Loginhandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/login"); err != nil {
			return
		}
		// this stops OPTIONS method from running and therefore not allowing to set access control allow headers
		// Prevents all request types other than POST
		// if r.Method != http.MethodPost {
		// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }

		if r.Method == http.MethodPost {
			fmt.Printf("----login-POST-----\n")
			var payload loginPayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				fmt.Println("decode prob login")
				log.Fatal(err)
			}
			fmt.Println(payload)

			email := payload.Email
			pw := payload.Pw

			fmt.Printf("Email: %s\n", email)
			fmt.Printf("password: %s\n", pw)
			type NotandAuthResp struct {
				Resp  AuthResponse  `json:"resp"`
				Notif []NotifStruct `json:"notif"`
			}
			var offlineNotResp NotandAuthResp
			var Resp AuthResponse
			Resp.Success = true
			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			var query *crud.Queries

			query = crud.New(db)

			// ### SEARCH DATABASE FROM USER ###

			curUser, err := query.GetUser(context.Background(), payload.Email)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to find user")
			}

			if curUser.Count < 1 {
				Resp.Success = false
				fmt.Println("Unable to find user")
			}

			// ### COMPARE PASSWORD WITH THE HASH IN THE DATABASE (SKIP IF USER NOT FOUND) ###

			err = bcrypt.CompareHashAndPassword([]byte(curUser.Password), []byte(payload.Pw))

			if err != nil {
				Resp.Success = false
				fmt.Println("Passwords do not match!")
			}

			Resp.UserId = int(curUser.ID)
			Resp.Fname = curUser.FirstName
			Resp.Lname = curUser.LastName
			Resp.Nname = curUser.NickName
			Resp.Avatar = curUser.Image
			Resp.Email = curUser.Email
			Resp.About = curUser.About
			Resp.Dob = curUser.Dob.String()
			Resp.Public = int(curUser.Public)

			if email == "f" {
				Resp.Success = false
			}

			// ### UPDATE SESSION COOKIE IN DATABASE AND BROWSER (SKIP IF USER NOT FOUND OR IF PASSWORD DOES NOT MATCH) ###
			sessionExist, err := query.SessionExists(context.Background(), curUser.ID)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to check session table!")
			}

			if Resp.Success {
				// add new session
				// create cookie
				var cookie SessionStruct

				cookie.SessionToken = uuid.NewV4().String()
				cookie.UserId = int(curUser.ID)

				if sessionExist > 0 {
					// update session in database
					var newSession crud.UpdateUserSessionParams
					newSession.UserID = int64(cookie.UserId)
					newSession.SessionToken = cookie.SessionToken
					query.UpdateUserSession(context.Background(), newSession)

				} else {
					// add session to database
					var session crud.CreateSessionParams
					fmt.Println("add session to database")
					session.SessionToken = cookie.SessionToken
					session.UserID = int64(cookie.UserId)
					_, err = query.CreateSession(context.Background(), session)
					if err != nil {
						fmt.Println("Unable to create session!")
					}
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    cookie.SessionToken,
					MaxAge:   34560000,
					Path:     "/",
					Secure:   true,
					SameSite: http.SameSiteNoneMode,
				})

			}

			offlineNotResp.Resp = Resp
			offlineNotResp.Notif = offlineNot(offlineNotResp.Resp.UserId)
			jsonResp, err := json.Marshal(offlineNotResp)
			fmt.Println(string(jsonResp))

			WriteHttpHeader(jsonResp, w)
		}
	}
}

func Reghandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/reg"); err != nil {
			return
		}

		// Prevents all request types other than POST
		// if r.Method != http.MethodPost {
		// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }

		if r.Method == http.MethodPost {
			fmt.Printf("----reg-POST-----\n")
			var payload regPayload
			var errormsg string = "error"
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(payload)

			email := payload.Email
			pw := payload.Pw
			fname := payload.Fname
			lname := payload.Lname
			dob := payload.Dob
			avatar := payload.Avatar
			nname := payload.Nname
			about := payload.About

			fmt.Printf("Email: %s\n", email)
			fmt.Printf("password: %s\n", pw)
			fmt.Printf("fname: %s\n", fname)
			fmt.Printf("lname: %s\n", lname)
			fmt.Printf("dob: %s\n", dob)
			fmt.Printf("avatar: %s\n", avatar)
			fmt.Printf("nname: %s\n", nname)
			fmt.Printf("about: %s\n", about)

			// used to run query
			var regPayload crud.CreateUserParams

			// will be used to respond
			var Resp AuthResponse
			Resp.Success = true
			// convert password using bcrypt
			password := []byte(payload.Pw)

			cryptPw, err := bcrypt.GenerateFromPassword(password, 10)

			if err != nil {
				Resp.Success = false
				errormsg = "ERROR - Invalid password"
				fmt.Println("Unable generate password!")
			}

			date, err := time.Parse("2006-01-02", payload.Dob)

			if err != nil {
				Resp.Success = false
				errormsg = "ERROR - Invalid date of birth"
				fmt.Println("Unable to convert date of birth")
			}

			regPayload.Password = string(cryptPw)
			regPayload.Email = payload.Email
			regPayload.FirstName = payload.Fname
			regPayload.LastName = payload.Lname
			regPayload.Dob = date
			regPayload.Image = payload.Avatar
			regPayload.NickName = payload.Nname
			regPayload.About = payload.About

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			var query *crud.Queries

			query = crud.New(db)

			// check if user already exists

			var checkExist crud.GetUserExistParams

			checkExist.Email = regPayload.Email
			checkExist.NickName = regPayload.NickName

			records, err := query.GetUserExist(context.Background(), checkExist)

			if err != nil {
				Resp.Success = false
				errormsg = "ERROR"
				fmt.Println("Unable to check if user exists because of error")
			}

			if records > 0 {
				errormsg = "ERROR - This email is already taken"
				fmt.Println("Unable to check if user exists because of duplicate")
				// user already exists
				Resp.Success = false

			} else {

				// ### ATTEMPT TO ADD USER TO DATABASE ###
				var curUser crud.User
				curUser, err := query.CreateUser(context.Background(), regPayload)

				if err != nil {
					Resp.Success = false
					errormsg = "ERROR - Something went wrong! Try again"
					fmt.Println("Unable to create user!")
				}

				Resp.UserId = int(curUser.ID)
				Resp.Fname = curUser.FirstName
				Resp.Lname = curUser.LastName
				Resp.Nname = curUser.NickName
				Resp.Avatar = curUser.Image
				Resp.Email = curUser.Email
				Resp.About = curUser.About
				Resp.Dob = curUser.Dob.String()

				if email == "f" {
					fmt.Println("email false")
					errormsg = "ERROR - Invalid e-mail"
					Resp.Success = false
				}
			}
			if Resp.Success == false {
				Resp.Fname = errormsg
			}
			jsonResp, err := json.Marshal(Resp)
			fmt.Println(string(jsonResp))

			WriteHttpHeader(jsonResp, w)
		}
	}
}

func Logouthandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)

		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/logout"); err != nil {
			return
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Prevents all request types other than POST
		// if r.Method != http.MethodGet {
		// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }

		// Declares the handler response
		Resp := AuthResponse{Success: true}

		c, err := r.Cookie("session_token")

		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionToken := c.Value

		// ### CONNECT TO DATABASE ###

		db := db.DbConnect()

		var query *crud.Queries

		query = crud.New(db)

		// ### REMOVE SESSION COOKIE FROM DATABASE AND BROWSER ###

		query.DeleteSession(context.Background(), sessionToken)

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: "",
		})

		// Marshals the response struct to a json object
		jsonResp, err := json.Marshal(Resp)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// Sets the http headers and writes the response to the browser
		WriteHttpHeader(jsonResp, w)
	}
}

func Posthandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		if r.Method == http.MethodPost {
			fmt.Printf("-----POST---(create-post)--\n")
			var payload PostStruct

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(payload)

			author := payload.Author
			message := payload.Message
			image := payload.Image
			privacy := payload.Privacy
			createdAt := time.Now()

			fmt.Printf("post author userid %d\n", author)
			fmt.Printf("post message %s\n", message)
			fmt.Printf("post image %s\n", image)
			fmt.Printf("post privacy %d\n", privacy)
			fmt.Printf("post created at %v\n", createdAt)

			var Resp PostResponse
			Resp.Success = true

			// insert post to database

			db := db.DbConnect()

			var post crud.CreatePostParams

			post.Author = int64(payload.Author)
			post.Message = payload.Message
			post.CreatedAt = createdAt
			post.Image = payload.Image
			post.Privacy = int64(payload.Privacy)

			query := crud.New(db)

			newPost, err := query.CreatePost(context.Background(), post)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to insert new post")
			}

			Resp.Author = int(newPost.Author)
			Resp.CreatedAt = newPost.CreatedAt.String()
			Resp.Image = newPost.Image
			Resp.Message = newPost.Message

			curUser, err := query.GetUserById(context.Background(), newPost.Author)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to get user information")
			}

			Resp.Avatar = curUser.Image
			Resp.Fname = curUser.FirstName
			Resp.Nname = curUser.NickName
			Resp.Lname = curUser.LastName

			jsonResp, err := json.Marshal(Resp)

			WriteHttpHeader(jsonResp, w)
		}

		if r.Method == http.MethodGet {
			fmt.Printf("----post-GET---(display-posts)--\n")
			userId := r.URL.Query().Get("id")
			fmt.Println("USERID: ", userId)
			if userId == "" {
				return
			}
			int_user_id, err := strconv.Atoi(userId)
			if err != nil {
				log.Fatal("error - user id ", err)
			}
			var data []PostResponse
			// get all public posts
			db := db.DbConnect()
			query := crud.New(db)
			posts, err := query.GetAllPosts(context.Background())
			fmt.Println("posts: ", posts)
			if err != nil {
				fmt.Println("Unable to get all posts")
			}
			for _, post := range posts {
				fmt.Println("msg: ", post.Message, "--bool: ", checkFollower(int(post.Author), int_user_id))
				if post.Privacy == 0 ||
					(post.Privacy == 1 && checkFollower(int(post.Author), int_user_id) > 0 ||
						post.Privacy == 2 && checkFollower(int(post.Author), int_user_id) == 2) {
					var newPost PostResponse
					newPost.Success = true
					newPost.Id = int(post.ID)
					newPost.Author = int(post.Author)
					newPost.Message = post.Message
					newPost.CreatedAt = post.CreatedAt.String()
					newPost.Image = post.Image
					newPost.Privacy = int(post.Privacy)

					curUser, err := query.GetUserById(context.Background(), post.Author)

					if err != nil {
						newPost.Success = false
						fmt.Println("Unable to get user information")
					}

					newPost.Avatar = curUser.Image
					newPost.Fname = curUser.FirstName
					newPost.Lname = curUser.LastName
					newPost.Nname = curUser.NickName

					data = append(data, newPost)
				}
			}

			// fmt.Printf("data %v\n", data)
			jsonResp, _ := json.Marshal(data)
			// fmt.Printf("posts resp %s\n", string(jsonResp))

			WriteHttpHeader(jsonResp, w)

		}
	}
}

func PostCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Post")
		// EnableCors(&w)
		if r.Method == http.MethodPost {
			// fmt.Printf("-----POST---(create-comment)--\n")
			var payload PostCommentStruct

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(payload)

			postid := payload.PostId
			userid := payload.UserId
			content := payload.Message
			image := payload.Image
			payload.CreatedAt = time.Now()

			fmt.Printf("postid %d\n", postid)
			fmt.Printf("userid %d\n", userid)
			fmt.Printf("content %s\n", content)
			fmt.Printf("image %s\n", image)

			// insert comment into database

			db := db.DbConnect()

			var postComment crud.CreatePostCommentParams

			postComment.PostID = int64(payload.PostId)
			postComment.UserID = int64(payload.UserId)
			postComment.Message = payload.Message
			postComment.CreatedAt = payload.CreatedAt
			postComment.Image = payload.Image

			query := crud.New(db)

			_, err = query.CreatePostComment(context.Background(), postComment)

			if err != nil {
				fmt.Println("Unable to insert new comment")
			}

			var Resp PostCommentResponse
			Resp.Success = true
			jsonResp, err := json.Marshal(Resp)

			WriteHttpHeader(jsonResp, w)
		}

		if r.Method == http.MethodGet {
			// fmt.Printf("----post-comment-GET---(display)--\n")

			var data []PostCommentResponse

			// get all comments

			db := db.DbConnect()

			query := crud.New(db)

			comments, err := query.GetAllComments(context.Background())

			if err != nil {
				fmt.Println("Unable to get all comments")
			}

			for _, comment := range comments {
				var newComment PostCommentResponse
				newComment.Success = true
				newComment.PostId = int(comment.PostID)
				newComment.UserId = int(comment.UserID)
				newComment.CreatedAt = comment.CreatedAt.String()
				newComment.Message = comment.Message
				newComment.Image = comment.Image
				newComment.Id = int(comment.ID)

				curUser, err := query.GetUserById(context.Background(), comment.UserID)

				if err != nil {
					newComment.Success = false
					fmt.Println("Unable to get user information")
				}

				newComment.Avatar = curUser.Image
				newComment.Fname = curUser.FirstName
				newComment.Lname = curUser.LastName
				newComment.Nname = curUser.NickName

				data = append(data, newComment)
			}
			jsonResp, _ := json.Marshal(data)
			// fmt.Printf("posts resp %s\n", string(jsonResp))

			WriteHttpHeader(jsonResp, w)
		}
	}
}

func Userhandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/user"); err != nil {
			return
		}

		// Prevents all request types other than GET
		if r.Method != http.MethodGet {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Checks to find a user id in the url
		fmt.Println("urluser: ", r.URL)
		userId := r.URL.Query().Get("id")
		id, err := strconv.Atoi(userId)
		if err != nil {
			fmt.Println("Unable to convert to int")
		}

		foundId := false

		if userId != "" {
			foundId = true
		}

		// Declares the payload struct
		var Resp UserPayload

		// ### CONNECT TO DATABASE ###

		db := db.DbConnect()

		query := crud.New(db)

		if foundId {
			// ### GET USER BY ID ###
			user, err := query.GetUserById(context.Background(), int64(id))

			if err != nil {
				fmt.Println("Unable to find user")
			}

			var oneUser UserStruct

			oneUser.Id = int(user.ID)
			oneUser.Fname = user.FirstName
			oneUser.Lname = user.LastName
			oneUser.Nname = user.NickName
			oneUser.Email = user.Email
			oneUser.Password = user.Password
			oneUser.Dob = user.Dob.String()
			oneUser.Avatar = user.Image
			oneUser.About = user.About
			oneUser.Public = int(user.Public)

			Resp.Data = append(Resp.Data, oneUser)

		} else {
			// ### GET ALL USERS ###
			users, err := query.ListUsers(context.Background())

			if err != nil {
				fmt.Println("Unable to get users")
			}

			for _, user := range users {
				var oneUser UserStruct

				oneUser.Id = int(user.ID)
				oneUser.Fname = user.FirstName
				oneUser.Lname = user.LastName
				oneUser.Nname = user.NickName
				oneUser.Email = user.Email
				oneUser.Password = user.Password
				oneUser.Dob = user.Dob.String()
				oneUser.Avatar = user.Image
				oneUser.About = user.About
				oneUser.Public = int(user.Public)

				Resp.Data = append(Resp.Data, oneUser)
			}

		}

		// Marshals the response struct to a json object
		jsonResp, err := json.Marshal(Resp)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// Sets the http headers and writes the response to the browser
		WriteHttpHeader(jsonResp, w)
	}
}

func UserFollowerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/user-follower"); err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:

			// Checks to find a user id in the url
			targetId := r.URL.Query().Get("id")
			id, err := strconv.Atoi(targetId)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}

			foundId := false

			if targetId != "" {
				foundId = true
			}

			// Declares the payload struct
			var Resp UserPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			if foundId {
				// ### GET USER FOLLOWERS ###
				followers, err := query.GetFollowers(context.Background(), int64(id))

				if err != nil {
					fmt.Println("Unable to find followers")
				}

				for _, follower := range followers {
					user, err := query.GetUserById(context.Background(), follower.SourceID)

					if err != nil {
						fmt.Println("Unable to find user")
					}
					if follower.Status > 0 {
						var oneUser UserStruct
						oneUser.Id = int(user.ID)
						oneUser.Fname = user.FirstName
						oneUser.Lname = user.LastName
						oneUser.Nname = user.NickName
						oneUser.Email = user.Email
						oneUser.Password = user.Password
						oneUser.Dob = user.Dob.String()
						oneUser.Avatar = user.Image
						oneUser.About = user.About
						oneUser.Public = int(user.Public)
						Resp.Data = append(Resp.Data, oneUser)
					}
				}
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the follower details and handler response
			var follower UserFollowerStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&follower)
			if err != nil {
				Resp.Success = false
			}

			fmt.Println("follow", follower)

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD FOLLOWER TO DATABASE ###
			var newFollower crud.CreateFollowerParams

			newFollower.SourceID = int64(follower.SourceId)
			newFollower.TargetID = int64(follower.TargetId)
			newFollower.Status = int64(follower.Status)
			// newFollower.LastMsgAt = follower.LastMsgAt

			_, err = query.CreateFollower(context.Background(), newFollower)

			if err != nil {
				fmt.Println("Unable to insert follower")
				Resp.Success = false
			}

		case http.MethodDelete:
			// Declares the variables to store the follower details and handler response
			var follower UserFollowerStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&follower)
			if err != nil {
				Resp.Success = false
			}

			fmt.Println("unfollow", follower)

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			var delFollower crud.DeleteFollowerParams

			delFollower.SourceID = int64(follower.SourceId)
			delFollower.TargetID = int64(follower.TargetId)

			err = query.DeleteFollower(context.Background(), delFollower)

			if err != nil {
				fmt.Println("Unable to delete follower")
				Resp.Success = false
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				fmt.Println("marshal")
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
			// default:
			// 	// Prevents all request types other than POST and GET
			// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// 	return
		}
	}
}

func UserFollowingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/user-following"); err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:

			// Checks to find a user id in the url
			sourceId := r.URL.Query().Get("id")
			id, err := strconv.Atoi(sourceId)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}

			foundId := false

			if sourceId != "" {
				foundId = true
			}

			// Declares the payload struct
			var Resp UserPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			if foundId {
				// ### GET USER FOLLOWERS ###
				followings, err := query.GetFollowings(context.Background(), int64(id))

				if err != nil {
					fmt.Println("Unable to find followers")
				}

				for _, following := range followings {
					user, err := query.GetUserById(context.Background(), following.TargetID)

					if err != nil {
						fmt.Println("Unable to find user")
					}
					if following.Status == 1 || following.Status == 2 {
						var oneUser UserStruct

						oneUser.Id = int(user.ID)
						oneUser.Fname = user.FirstName
						oneUser.Lname = user.LastName
						oneUser.Nname = user.NickName
						oneUser.Email = user.Email
						oneUser.Password = user.Password
						oneUser.Dob = user.Dob.String()
						oneUser.Avatar = user.Image
						oneUser.About = user.About
						oneUser.Public = int(user.Public)

						Resp.Data = append(Resp.Data, oneUser)
					}
				}

			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the follower details and handler response
			var follower UserFollowerStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&follower)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### delete FOLLOWER TO DATABASE ###

			var newFollower crud.DeleteFollowerParams

			newFollower.SourceID = int64(follower.SourceId)
			newFollower.TargetID = int64(follower.TargetId)

			err = query.DeleteFollower(context.Background(), newFollower)
			fmt.Println("NEW FOLLOW REQUESTED")
			if err != nil {
				fmt.Println("Unable to insert follower")
				Resp.Success = false
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func CloseFriendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/close-friend"); err != nil {
			return
		}

		fmt.Println(r.Method)
		switch r.Method {
		case http.MethodGet:

			// Checks to find a user id in the url
			sourceId := r.URL.Query().Get("id")
			id, err := strconv.Atoi(sourceId)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}
			foundId := false

			if sourceId != "" {
				foundId = true
			}

			// Declares the payload struct
			var Resp UserPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			if foundId {
				// ### GET USER FOLLOWERS ###
				followings, err := query.GetFollowings(context.Background(), int64(id))

				if err != nil {
					fmt.Println("Unable to find followers")
				}

				for _, following := range followings {
					user, err := query.GetUserById(context.Background(), following.TargetID)

					if err != nil {
						fmt.Println("Unable to find user")
					}
					if following.Status == 2 {
						var oneUser UserStruct

						oneUser.Id = int(user.ID)
						oneUser.Fname = user.FirstName
						oneUser.Lname = user.LastName
						oneUser.Nname = user.NickName
						oneUser.Email = user.Email
						oneUser.Password = user.Password
						oneUser.Dob = user.Dob.String()
						oneUser.Avatar = user.Image
						oneUser.About = user.About
						oneUser.Public = int(user.Public)

						Resp.Data = append(Resp.Data, oneUser)
					}
				}

			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			type Ids struct {
				SourceID int64 `json:"sourceid"`
				TargetID int64 `json:"targetid"`
			}

			// Declares the variables to store the follower details and handler response
			var followers Ids
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&followers)
			if err != nil {
				Resp.Success = false
			}
			fmt.Println("closefriend: ", followers)

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			followerContext := crud.CheckFollowerParams{SourceID: followers.SourceID, TargetID: followers.TargetID}
			follower, err := query.CheckFollower(context.Background(), followerContext)
			if err != nil {
				fmt.Println("Unable to find follower")
				Resp.Success = false
			}

			// ### update FOLLOWER TO DATABASE ###

			var newFollower crud.UpdateFollowerParams
			newFollower.SourceID = follower.SourceID
			newFollower.TargetID = follower.TargetID

			if follower.Status == 1 {
				newFollower.Status = int64(2)
			} else if follower.Status == 2 {
				newFollower.Status = int64(1)
			} else {
				newFollower.Status = int64(follower.Status)
			}

			_, err = query.UpdateFollower(context.Background(), newFollower)
			fmt.Println("UPDATED CLOSE FRIEND")
			if err != nil {
				fmt.Println("Unable to update close friend")
				Resp.Success = false
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
			// default:
			// 	// Prevents all request types other than POST and GET
			// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// 	return
		}
	}
}

func UserMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/user-message"); err != nil {
			return
		}

		targetId := r.URL.Query().Get("targetid")
		if targetId == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		tId, err := strconv.Atoi(targetId)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		sourceId := r.URL.Query().Get("sourceid")
		if sourceId == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		sId, err := strconv.Atoi(sourceId)
		if err != nil {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			// Declares the payload struct
			var allMessages []crud.UserMessage

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET ALL MESSAGES FOR THE TARGET ID AND SOURCE ID ####
			var msg crud.GetMessagesParams

			msg.SourceID = int64(sId)
			msg.SourceID_2 = int64(tId)
			msg.TargetID = int64(tId)
			msg.TargetID_2 = int64(sId)

			allMessages, err = query.GetMessages(context.Background(), msg)

			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			var messages UserMessagePayload

			for _, message := range allMessages {
				var newMessage UserMessageStruct

				newMessage.Id = int(message.ID)
				newMessage.TargetId = int(message.TargetID)
				newMessage.SourceId = int(message.SourceID)
				newMessage.Message = message.Message
				newMessage.CreatedAt = message.CreatedAt.String()

				messages.Data = append(messages.Data, newMessage)
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(messages)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the user message details and handler response
			var userMessage UserMessageStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&userMessage)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD USER MESSAGE TO DATABASE ###

			var message crud.CreateMessageParams
			message.CreatedAt = time.Now()
			message.Message = userMessage.Message
			message.SourceID = int64(userMessage.SourceId)
			message.TargetID = int64(userMessage.TargetId)

			_, err = query.CreateMessage(context.Background(), message)

			if err != nil {
				Resp.Success = false
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func Grouphandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group"); err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Checks to find a group id in the url
			groupId := r.URL.Query().Get("id")
			foundId := false

			if groupId != "" {
				foundId = true

			}

			gId, err := strconv.Atoi(groupId)

			if err != nil {
				fmt.Println("Unable to convert group ID")
			}

			var Resp GroupPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// Gets the group by id if an id was passed in the url
			// Otherwise, gets all group
			if foundId {
				// Declares the payload struct
				var group crud.Group
				// GET USER ID FROM SESSION

				session, err := r.Cookie("session_token")
				if err != nil {
					fmt.Println(err)
					return
				}

				sessionTable, err := query.GetUserId(context.Background(), session.Value)

				if err != nil {
					http.Error(w, "500 internal server error", http.StatusInternalServerError)
					return
				}

				// ### CHECK IF USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###
				var member crud.CheckIfMemberParams

				member.GroupID = int64(gId)
				member.Status = 1
				member.UserID = sessionTable.UserID

				groupData, err := query.CheckIfMember(context.Background(), member)

				fmt.Println("-------", member)

				if err != nil {
					fmt.Println("Unable to get group data")
				}

				// ### IF THEY MATCH, GET GROUP DATA FROM DATABASE ###

				if groupData == 1 {
					// ### GET GROUP BY ID ###
					group, err = query.GetGroup(context.Background(), int64(gId))

					if err != nil {
						fmt.Println("Unable to get group id")
					}

					var newGroup GroupStruct

					newGroup.Id = int(group.ID)
					newGroup.Title = group.Title
					newGroup.Creator = int(group.Creator)
					newGroup.Description = group.Description
					newGroup.CreatedAt = group.CreatedAt.String()

					Resp.Data = append(Resp.Data, newGroup)

				} else {
					// ### ELSE, REQUEST TO JOIN ###
					//empty response
				}

			} else {
				// ### GET ALL GROUPS ###

				groups, err := query.GetAllGroups(context.Background())

				if err != nil {
					fmt.Println("Unable to get groups")
				}

				for _, group := range groups {
					var newGroup GroupStruct

					newGroup.Id = int(group.ID)
					newGroup.Title = group.Title
					newGroup.Creator = int(group.Creator)
					newGroup.Description = group.Description
					newGroup.CreatedAt = group.CreatedAt.String()

					Resp.Data = append(Resp.Data, newGroup)
				}
				fmt.Println(Resp)
			}

			jsonResp, err := json.Marshal(Resp)

			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			WriteHttpHeader(jsonResp, w)

		case http.MethodPost:
			// Declares the variables to store the group details and handler response
			fmt.Println("Post Grp")

			var group GroupStruct
			Resp := GroupResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			fmt.Printf("grp body %v\n", r.Body)
			err := json.NewDecoder(r.Body).Decode(&group)
			if err != nil {
				fmt.Println("decode prob grp") // why? but it can be stored in db
				// Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP TO DATABASE ###

			var groupData crud.CreateGroupParams

			groupData.CreatedAt = time.Now()
			groupData.Creator = int64(group.Creator)
			groupData.Description = group.Description
			groupData.Title = group.Title

			newGroup, err := query.CreateGroup(context.Background(), groupData)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to create new group")
			}

			// ### ADD GROUP CREATOR TO GROUP MEMBER TABLE ###

			var creator crud.CreateGroupMemberParams

			creator.GroupID = newGroup.ID
			creator.Status = 1
			creator.UserID = newGroup.Creator

			Resp.CreatedGroupId = int(newGroup.ID)
			Resp.Creator = int(newGroup.Creator)

			_, err = query.CreateGroupMember(context.Background(), creator)

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to add creator to members list")
			}
			fmt.Printf("resp %v\n", Resp)
			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Printf("jsonresp %v", jsonResp)
			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// return
		}
	}
}

func GroupMemberHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-member"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		switch r.Method {
		case http.MethodGet:
			fmt.Println("url:\n", r.URL)
			// Checks to find a user id in the url
			userId := r.URL.Query().Get("userid")
			groupId := r.URL.Query().Get("groupid")
			uId, err := strconv.Atoi(userId)

			if err != nil {
				fmt.Println("Unable to convert user ID")
			}

			gId, err := strconv.Atoi(groupId)

			if err != nil {
				fmt.Println("Unable to convert group ID")
			}

			foundUserId := false
			foundGroupId := false
			if userId != "" {
				foundUserId = true
			}

			if groupId != "" {
				foundGroupId = true
			}
			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// gets all groups user is a member of
			if foundUserId {
				groups, err := query.GetAllGroupsByUser(context.Background(), int64(uId))
				if err != nil {
					fmt.Println("Unable to get groups")
				}

				var groupsResp GroupPayload

				for _, group := range groups {
					var oneGroup GroupStruct

					oneGroup.CreatedAt = group.CreatedAt.String()
					oneGroup.Description = group.Description
					oneGroup.Creator = int(group.Creator)
					oneGroup.Id = int(group.ID)
					oneGroup.Title = group.Title

					groupsResp.Data = append(groupsResp.Data, oneGroup)
				}

				// Marshals the response struct to a json object
				jsonResp, err := json.Marshal(groupsResp)
				if err != nil {
					http.Error(w, "500 internal server error", http.StatusInternalServerError)
					return
				}

				// Sets the http headers and writes the response to the browser
				WriteHttpHeader(jsonResp, w)

			}

			// get all members with the following group id
			if foundGroupId {
				users, err := query.GetGroupMembersByGroupId(context.Background(), crud.GetGroupMembersByGroupIdParams{
					GroupID: int64(gId),
					Status:  1,
				})

				if err != nil {
					fmt.Println("Unable to get members")
				}

				// Declares the payload struct
				type UserPayload2 struct {
					Members    []UserStruct `json:"members"`
					NotMembers []UserStruct `json:"notmembers"`
				}
				var usersResp UserPayload2

				for _, user := range users {
					var oneUser UserStruct

					oneUser.Id = int(user.ID)
					oneUser.Fname = user.FirstName
					oneUser.Lname = user.LastName
					oneUser.Nname = user.NickName
					oneUser.Email = user.Email
					oneUser.Password = user.Password
					oneUser.Dob = user.Dob.String()
					oneUser.Avatar = user.Image
					oneUser.About = user.About
					oneUser.Public = int(user.Public)

					usersResp.Members = append(usersResp.Members, oneUser)
				}
				users2, err := query.GetGroupMembersByGroupIdWithoutStatus(context.Background(), int64(gId))

				if err != nil {
					fmt.Println("Unable to get members")
				}
				allUsers, err := query.ListUsers(context.Background())
				for i := 0; i < len(allUsers); i++ {
					check := false
					for k := 0; k < len(users2); k++ {
						if allUsers[i].ID == users2[k].ID {
							check = true
						} else if k == len(users2)-1 && !check {
							var oneUser UserStruct

							oneUser.Id = int(allUsers[i].ID)
							oneUser.Fname = allUsers[i].FirstName
							oneUser.Lname = allUsers[i].LastName
							oneUser.Nname = allUsers[i].NickName
							oneUser.Email = allUsers[i].Email
							oneUser.Password = allUsers[i].Password
							oneUser.Dob = allUsers[i].Dob.String()
							oneUser.Avatar = allUsers[i].Image
							oneUser.About = allUsers[i].About
							oneUser.Public = int(allUsers[i].Public)

							usersResp.NotMembers = append(usersResp.NotMembers, oneUser)
						}
					}
				}
				// Marshals the response struct to a json object
				jsonResp, err := json.Marshal(usersResp)
				if err != nil {
					http.Error(w, "500 internal server error", http.StatusInternalServerError)
					return
				}

				// Sets the http headers and writes the response to the browser
				WriteHttpHeader(jsonResp, w)
			}

			// case http.MethodPost:
			// 	// Declares the variables to store the group member details and handler response
			// 	var groupMember GroupMemberStruct
			// 	Resp := AuthResponse{Success: true}
			// 	// Decodes the json object to the struct, changing the response to false if it fails
			// 	err := json.NewDecoder(r.Body).Decode(&groupMember)
			// 	if err != nil {
			// 		Resp.Success = false
			// 	}

			// 	// ### CONNECT TO DATABASE ###

			// 	// db := db.DbConnect()

			// 	// query := crud.New(db)
			// 	// statusStr := strconv.Itoa(groupMember.Status)
			// 	// ### UPDATE GROUP REQUEST TABLE AND ADD USER TO GROUP MEMBER TABLE ###

			// 	// _, err = query.UpdateGroupRequest(context.Background(), crud.UpdateGroupRequestParams{
			// 	// 	Status:  statusStr,
			// 	// 	GroupID: int64(groupMember.GroupId),
			// 	// 	UserID:  int64(groupMember.UserId),
			// 	// })

			// 	// if err != nil {
			// 	// 	Resp.Success = false
			// 	// 	fmt.Println("Unable to update group request")
			// 	// }

			// 	// _, err = query.CreateGroupMember(context.Background(), crud.CreateGroupMemberParams{
			// 	// 	UserID:  int64(groupMember.UserId),
			// 	// 	GroupID: int64(groupMember.GroupId),
			// 	// 	Status:  1,
			// 	// })

			// 	// if err != nil {
			// 	// 	Resp.Success = false
			// 	// 	fmt.Println("Unable to insert group member")
			// 	// }

			// 	// Marshals the response struct to a json object
			// 	jsonResp, err := json.Marshal(Resp)
			// 	if err != nil {
			// 		http.Error(w, "500 internal server error", http.StatusInternalServerError)
			// 		return
			// 	}

			// 	// Sets the http headers and writes the response to the browser
			// 	WriteHttpHeader(jsonResp, w)
			// default:
			// 	// Prevents all request types other than POST and GET
			// 	// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// 	// return
		}
	}
}

func GroupRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-request"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		switch r.Method {
		case http.MethodGet:
			// ### CHECK USER IS GROUP CREATOR ###

			// get group id from url
			groupId := r.URL.Query().Get("groupid")

			gId, err := strconv.Atoi(groupId)

			if err != nil {
				fmt.Println("Unable to convert group ID")

			}

			// connect to database

			db := db.DbConnect()

			query := crud.New(db)

			// get user from cookie

			session, err := r.Cookie("SessionToken")

			sessionTable, err := query.GetUserId(context.Background(), session.Value)

			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			uId := sessionTable.UserID

			creator, err := query.CheckIfCreator(context.Background(), crud.CheckIfCreatorParams{
				Creator: uId,
				ID:      int64(gId),
			})

			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Declares the payload struct
			var Resp GroupRequestPayload

			// ### GET ALL GROUP REQUESTS FOR GROUP ID ###
			if creator > 0 {
				groups, err := query.GetAllGroupRequests(context.Background(), int64(gId))

				if err != nil {
					fmt.Println("Unable to get groups")
				}

				for _, group := range groups {
					var oneGroup GroupRequestStruct

					oneGroup.Id = int(group.ID)
					oneGroup.UserId = int(group.UserID)
					oneGroup.GroupId = int(group.GroupID)
					oneGroup.Status = group.Status

					Resp.Data = append(Resp.Data, oneGroup)
				}

			}
			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the group request details and handler response
			var groupRequest GroupRequestStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupRequest)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP REQUEST TO DATABASE ###

			_, err = query.CreateGroupRequest(context.Background(), crud.CreateGroupRequestParams{
				UserID:  int64(groupRequest.UserId),
				GroupID: int64(groupRequest.GroupId),
				Status:  groupRequest.Status,
			})

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to add new group request")
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// return
		}
	}
}

func GroupPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-post"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		switch r.Method {
		case http.MethodGet:
			// Checks to find a group id in the url
			groupId := r.URL.Query().Get("groupid")
			postId := r.URL.Query().Get("postid")

			gId, err := strconv.Atoi(groupId)

			if err != nil {
				fmt.Println("Unable to convert group ID")
			}

			pId, err := strconv.Atoi(postId)

			if err != nil {
				fmt.Println("Unable to convert post ID")
			}

			foundId := false

			if postId != "" {
				foundId = true
			}

			// Declares the payload struct
			var Resp GroupPostPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// Gets the post by id if an id was passed in the url

			// Otherwise, gets all posts
			if foundId {
				// ### GET GROUP POST BY ID ###
				groupPost, err := query.GetGroupPostById(context.Background(), int64(pId))

				if err != nil {
					fmt.Println("Unable to get group post")
				}

				var onePost GroupPostStruct
				us, err := query.GetUserById(context.Background(), groupPost.Author)
				if err != nil {
					log.Fatal(err)
				}
				onePost.Fname = us.FirstName
				onePost.Lname = us.LastName
				onePost.Nickname = us.NickName
				onePost.Id = int(groupPost.ID)
				onePost.GroupId = int(groupPost.GroupID)
				onePost.Author = int(groupPost.Author)
				onePost.Message = groupPost.Message
				onePost.Image = groupPost.Image
				onePost.CreatedAt = groupPost.CreatedAt.String()

				Resp.Data = append(Resp.Data, onePost)

			} else {
				// ### GET ALL GROUP POSTS ###
				groupPosts, err := query.GetGroupPosts(context.Background(), int64(gId))

				if err != nil {
					fmt.Println("Unable to get group post")
				}

				for _, post := range groupPosts {
					var onePost GroupPostStruct
					us, err := query.GetUserById(context.Background(), post.Author)
					if err != nil {
						log.Fatal(err)
					}
					onePost.Fname = us.FirstName
					onePost.Lname = us.LastName
					onePost.Nickname = us.NickName
					onePost.Id = int(post.ID)
					onePost.GroupId = int(post.GroupID)
					onePost.Author = int(post.Author)
					onePost.Message = post.Message
					onePost.Image = post.Image
					onePost.CreatedAt = post.CreatedAt.String()

					Resp.Data = append(Resp.Data, onePost)
				}

			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the group post details and handler response
			var groupPost GroupPostStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupPost)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP POST TO DATABASE ###

			_, err = query.CreateGroupPost(context.Background(), crud.CreateGroupPostParams{
				Author:    int64(groupPost.Author),
				GroupID:   int64(groupPost.GroupId),
				Message:   groupPost.Message,
				Image:     groupPost.Image,
				CreatedAt: time.Now(),
			})

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to create group post")
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// return
		}
	}
}

func GroupPostCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-post-comment"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		// Checks to find a post id in the url

		switch r.Method {
		case http.MethodGet:
			groupPostId := r.URL.Query().Get("id")
			gPostId, err := strconv.Atoi(groupPostId)

			if err != nil {
				fmt.Println("Unable to convert group post ID")
			}

			if groupPostId == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}

			// Declares the payload struct
			var Resp GroupPostCommentPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET ALL COMMENTS FOR THE GROUP POST ID ###

			comments, err := query.GetGroupPostComments(context.Background(), int64(gPostId))

			if err != nil {
				fmt.Println("Unable to get comments")
			}

			for _, comment := range comments {
				us, err := query.GetUserById(context.Background(), comment.Author)
				if err != nil {
					log.Fatal(err)
				}
				var newComment GroupPostCommentStruct
				newComment.Fname = us.FirstName
				newComment.Lname = us.LastName
				newComment.Nickname = us.NickName
				newComment.Id = int(comment.ID)
				newComment.Author = int(comment.Author)
				newComment.Message = comment.Message
				newComment.GroupPostId = int(comment.GroupPostID)
				newComment.CreatedAt = comment.CreatedAt.String()
				newComment.AuthorAvatar = us.Image
				Resp.Data = append(Resp.Data, newComment)
			}
			fmt.Println("Resp: ", Resp.Data, len(Resp.Data))
			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the group post comment details and handler response
			var groupPostComment GroupPostCommentStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupPostComment)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP POST COMMENT TO DATABASE ###

			_, err = query.CreateGroupPostComment(context.Background(), crud.CreateGroupPostCommentParams{
				Author:      int64(groupPostComment.Author),
				GroupPostID: int64(groupPostComment.GroupPostId),
				Message:     groupPostComment.Message,
				CreatedAt:   time.Now(),
			})

			if err != nil {
				fmt.Println("Unable to create comment")
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GroupEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-event"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		switch r.Method {
		case http.MethodGet:
			// Declares the payload struct
			groupId := r.URL.Query().Get("id")
			if groupId == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}
			gId, err := strconv.Atoi(groupId)

			userid := r.URL.Query().Get("userid")
			if userid == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}
			uid, err := strconv.Atoi(userid)

			if err != nil {
				fmt.Println("Unable to convert group ID")
			}

			var Resp GroupEventPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET ALL EVENTS FOR THE GROUP ID ###

			events, err := query.GetGroupEvents(context.Background(), int64(gId))
			for _, event := range events {
				var newEvent GroupEventStruct
				err = query.ExecUpdateGroupEventMember(context.Background(), crud.ExecUpdateGroupEventMemberParams{UserID: int64(uid), EventID: event.ID})
				newEvent.Id = int(event.ID)
				newEvent.GroupId = int(event.GroupID)
				newEvent.Author = int(event.Author)
				newEvent.Title = event.Title
				newEvent.Description = event.Description
				newEvent.CreatedAt = event.CreatedAt.String()
				newEvent.Date = event.Date.String()

				Resp.Data = append(Resp.Data, newEvent)

			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the group event details and handler response
			var groupEvent GroupEventStruct
			Resp := AuthResponse{Success: true}
			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupEvent)
			if err != nil {
				Resp.Success = false
			}
			date_int, err := strconv.Atoi(groupEvent.CreatedAt)
			if err != nil {
				log.Fatal(err)
			}
			jsDateMs := int64(date_int)
			goDate := time.Unix(0, jsDateMs*int64(time.Millisecond))
			fmt.Println("CREATED AT on handler func: ", goDate)
			newdate := strings.Split(groupEvent.Date, "T")
			dateNew := newdate[0] + " " + newdate[1] + ":00"
			fmt.Println(dateNew)
			date, err := time.Parse("2006-01-02 15:04:05", dateNew)
			fmt.Println(date)
			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to convert date")
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP EVENT TO DATABASE ###

			_, err = query.CreateGroupEvent(context.Background(), crud.CreateGroupEventParams{
				Author:      int64(groupEvent.Author),
				GroupID:     int64(groupEvent.GroupId),
				Title:       groupEvent.Title,
				Description: groupEvent.Description,
				CreatedAt:   goDate,
				Date:        date,
			})

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to create event")
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// return
		}
	}
}

func GroupEventMemberHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-event-member"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###
		switch r.Method {
		case http.MethodGet:
			fmt.Println("url:=", r.URL)
			var checkforevent bool = true
			// Checks to find a post id in the url
			eventId := r.URL.Query().Get("id")
			fmt.Println("bool:= ", eventId == "", " len:= ", len(eventId))
			if eventId == "" {
				fmt.Println("checking events by user")
				checkforevent = false
			}
			eId, err := strconv.Atoi(eventId)
			userId := r.URL.Query().Get("userid")
			uId, err := strconv.Atoi(userId)
			if err != nil {
			}
			var Resp GroupEventMemberPayload
			var Resp2 GroupEventPayload
			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET ALL GROUP EVENT MEMBERS ###
			if checkforevent {
				members, err := query.GetGroupEventMembers(context.Background(), int64(eId))

				if err != nil {
					fmt.Println("Unable to get all members")
				}

				for _, member := range members {
					var newMember GroupEventMemberStruct

					newMember.Id = int(member.ID)
					newMember.Status = int(member.Status)
					newMember.UserId = int(member.UserID)
					newMember.EventId = int(member.EventID)

					Resp.Data = append(Resp.Data, newMember)

				}
				// Marshals the response struct to a json object
				jsonResp, err := json.Marshal(Resp)
				if err != nil {
					http.Error(w, "500 internal server error", http.StatusInternalServerError)
					return
				}
				WriteHttpHeader(jsonResp, w)

			} else {
				acceptedEventsByUser, err := query.GetGroupEventsByUserAccepted(context.Background(), int64(uId))
				if err != nil {
					log.Fatal(err)
				}
				for _, oneeventId := range acceptedEventsByUser {
					selectedEvent, err := query.GetGroupEventById(context.Background(), oneeventId.EventID)
					if err != nil {
					}
					// s2, _ := json.MarshalIndent(selectedEvent, "", "\t")
					// fmt.Print("oneevent:=  ", string(s2))
					var oneevent GroupEventStruct
					oneevent.Id = int(oneeventId.EventID)
					oneevent.Title = selectedEvent.Title
					oneevent.Description = selectedEvent.Description
					oneevent.Date = selectedEvent.Date.String()
					oneevent.GroupId = int(selectedEvent.GroupID)
					// s, _ := json.MarshalIndent(oneevent, "", "\t")
					// fmt.Print("oneevent:=  ", string(s))
					Resp2.Data = append(Resp2.Data, oneevent)

				}

				// Marshals the response struct to a json object
				jsonResp, err := json.Marshal(Resp2)
				if err != nil {
					http.Error(w, "500 internal server error", http.StatusInternalServerError)
					return
				}
				WriteHttpHeader(jsonResp, w)

			}

			// Sets the http headers and writes the response to the browser

		case http.MethodPost:
			// Declares the variables to store the group event member details and handler response
			var groupEventMember GroupEventMemberStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupEventMember)
			if err != nil {
				log.Fatal(err)
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD/UPDATE GROUP EVENT MEMBER TO DATABASE ###

			exists, err := query.GetGroupEventMember(context.Background(), crud.GetGroupEventMemberParams{
				EventID: int64(groupEventMember.EventId),
				UserID:  int64(groupEventMember.UserId),
			})

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to get event member")
			}

			if exists > 0 {
				// update
				_, err = query.UpdateGroupEventMember(context.Background(), crud.UpdateGroupEventMemberParams{
					Status:  int64(groupEventMember.Status),
					EventID: int64(groupEventMember.EventId),
					UserID:  int64(groupEventMember.UserId),
				})

				if err != nil {
					Resp.Success = false
					fmt.Println("Unable to update event member")
				}

			} else {
				// add
				_, err = query.CreateGroupEventMember(context.Background(), crud.CreateGroupEventMemberParams{
					Status:  int64(groupEventMember.Status),
					EventID: int64(groupEventMember.EventId),
					UserID:  int64(groupEventMember.UserId),
				})

				if err != nil {
					Resp.Success = false
					fmt.Println("Unable to add event member")
				}

			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			// http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// return
		}
	}
}

func GroupMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-message"); err != nil {
			return
		}

		// ### CHECK USER ID AND GROUP ID MATCH IN GROUP MEMBER TABLE ###

		groupId := r.URL.Query().Get("id")
		if groupId == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		gId, err := strconv.Atoi(groupId)

		if err != nil {
			fmt.Println("Unable to convert group ID")
		}

		switch r.Method {
		case http.MethodGet:
			// Declares the payload struct
			var Resp GroupMessagePayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### GET ALL MESSAGES FOR THE GROUP ID ###

			messages, err := query.GetGroupMessages(context.Background(), int64(gId))

			if err != nil {
				fmt.Println("Unable to get group messages")
			}

			for _, message := range messages {
				var newMessage GroupMessageStruct

				newMessage.Id = int(message.ID)
				newMessage.Message = message.Message
				newMessage.SourceId = int(message.SourceID)
				newMessage.GroupId = int(message.GroupID)
				newMessage.CreatedAt = message.CreatedAt.String()

				Resp.Data = append(Resp.Data, newMessage)
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		case http.MethodPost:
			// Declares the variables to store the group message details and handler response
			var groupMessage GroupMessageStruct
			Resp := AuthResponse{Success: true}

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupMessage)
			if err != nil {
				Resp.Success = false
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### ADD GROUP MESSAGE TO DATABASE ###

			_, err = query.CreateGroupMessage(context.Background(), crud.CreateGroupMessageParams{
				SourceID:  int64(groupMessage.SourceId),
				GroupID:   int64(groupMessage.GroupId),
				Message:   groupMessage.Message,
				CreatedAt: time.Now(),
			})

			if err != nil {
				Resp.Success = false
				fmt.Println("Unable to create group message")
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
func PrivacyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		fmt.Println(r.Method)
		if r.Method == "POST" {
			fmt.Println("its arrived")

			// Prevents the endpoint being called from other url paths
			if err := UrlPathMatcher(w, r, "/privacy"); err != nil {
				return
			}

			// // Prevents all request types other than GET
			// if r.Method != http.MethodPost {
			// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// 	return
			// }

			var payload UserId
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {

			}
			fmt.Println(payload)

			id := payload.TargetId
			public := payload.Public

			// Declares the payload struct
			var Resp AuthResponse

			Resp.Success = true

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// update public column in user table

			user, err := query.UpdateUserPrivacy(context.Background(), crud.UpdateUserPrivacyParams{
				Public: int64(public),
				ID:     int64(id),
			})

			if err != nil {
				fmt.Println("Unable to update user")
				Resp.Success = false
			}

			Resp.UserId = int(user.ID)
			Resp.Fname = user.FirstName
			Resp.Lname = user.LastName
			Resp.Nname = user.NickName
			Resp.Email = user.Email
			Resp.Public = int(user.Public)

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		}
	}
}
func UserFollowerStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/user-follow-status"); err != nil {
			return
		}
		if r.Method == "GET" {

			// Checks to find a user id in the url
			tid := r.URL.Query().Get("tid")
			sid := r.URL.Query().Get("sid")
			id, err := strconv.Atoi(tid)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}
			source_id, err := strconv.Atoi(sid)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}

			foundId := false

			if tid != "" {
				foundId = true
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			if foundId {
				// ### GET USER FOLLOWERS ###
				followers, err := query.GetFollowers(context.Background(), int64(id))
				var value string
				if err != nil {
					fmt.Println("Unable to find followers")
				}
				for _, follower := range followers {
					if err != nil {
						fmt.Println("Unable to find user")
					}
					// w.Header().Set("Content-Type", "application/json")
					if int(follower.SourceID) == source_id {
						if follower.Status == 1 {
							value = "false"
						} else if follower.Status == 0 {
							value = "true"
						}
					}
				}
				followers2, err := query.GetFollowers(context.Background(), int64(source_id))
				if err != nil {
					fmt.Println("Unable to find followers")
				}
				for _, follower := range followers2 {
					if err != nil {
						fmt.Println("Unable to find user")
					}
					if int(follower.SourceID) == id {
						if follower.Status == 2 {
							value += "closefriend"
						}

					}
				}
				fmt.Fprint(w, value)

			}
		}
	}
}
func checkFollower(targetid, sourceid int) int {
	fmt.Printf("checking target: %d and source: %d\n", targetid, sourceid)
	if sourceid == targetid {
		fmt.Println("OWN PROFILE")
		return 2
	}
	db := db.DbConnect()
	var query *crud.Queries
	query = crud.New(db)
	var checkFollower crud.CheckFollowerParams
	checkFollower.SourceID = int64(sourceid)
	checkFollower.TargetID = int64(targetid)
	following, err := query.CheckFollower(context.Background(), checkFollower)
	if err != nil {
		return 0
	}
	s, _ := json.MarshalIndent(following, "", "\t")
	fmt.Print("following: ", string(s))
	if following.ID != 0 {
		return int(following.Status)
	}
	return 0
}

func GroupRequestByUserHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-request-by-user"); err != nil {
			os.Exit(1)
			return
		}
		if r.Method == "GET" {

			userId := r.URL.Query().Get("id")
			if userId == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}
			uId, err := strconv.Atoi(userId)
			if err != nil {
				fmt.Println("Unable to convert group ID", uId)
			}
			db := db.DbConnect()

			query := crud.New(db)
			var Resp GroupRequestPayload
			groupRequests, err := query.GetGroupRequestsByUser(context.Background(), int64(uId))
			for _, group := range groupRequests {
				var newGroup GroupRequestStruct

				newGroup.Id = int(group.ID)
				newGroup.UserId = int(group.UserID)
				newGroup.GroupId = int(group.GroupID)
				newGroup.Status = "0"
				newGroup.CreatedAt = ""
				Resp.Data = append(Resp.Data, newGroup)
			}
			fmt.Println("---------", Resp.Data)
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		}
	}
}
func offlineNot(userid int) []NotifStruct {
	var ResultArr []NotifStruct
	db := db.DbConnect()

	query := crud.New(db)
	requests, err := query.GetAllGroupReq(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	groups, err := query.GetAllGroups(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, onegroup := range groups {
		if onegroup.Creator == int64(userid) {
			for _, onereq := range requests {
				if onereq.GroupID == onegroup.ID {
					var oneNotif NotifStruct
					oneNotif.Label = "noti"
					oneNotif.Id = 0
					oneNotif.Type = "join-req"
					oneNotif.TargetId = userid
					oneNotif.SourceId = int(onereq.UserID)
					oneNotif.Accepted = false
					oneNotif.CreatedAt = "not now"
					oneNotif.GroupId = int(onereq.GroupID)
					ResultArr = append(ResultArr, oneNotif)
				}
			}
		}
	}
	events, err := query.GetGroupEventsByUserNoReply(context.Background(), int64(userid))
	if err != nil {
		log.Fatal(err)
	}
	for _, oneevent := range events {
		groupid, err := query.GetGroupEventById(context.Background(), oneevent.EventID)
		if err != nil {
			log.Fatal(err)
		}
		if oneevent.Status == 0 {
			var oneNotif NotifStruct
			oneNotif.Label = "noti"
			oneNotif.Id = 0
			oneNotif.Type = "event-notif+" + groupid.Title
			oneNotif.TargetId = 987
			oneNotif.SourceId = userid
			oneNotif.Accepted = false
			oneNotif.CreatedAt = "not now"
			oneNotif.GroupId = int(groupid.GroupID)
			ResultArr = append(ResultArr, oneNotif)
		}
	}
	invites, err := query.GetGroupMembersByUserId(context.Background(), crud.GetGroupMembersByUserIdParams{UserID: int64(userid), Status: 0})
	if err != nil {
		log.Fatal(err)
	}
	for _, oneinvite := range invites {
		group, err := query.GetGroup(context.Background(), oneinvite.GroupID)
		if err != nil {
			log.Fatal(err)
		}
		var oneNotif NotifStruct
		oneNotif.Label = "noti"
		oneNotif.Id = 0
		oneNotif.Type = "invitation"
		oneNotif.TargetId = userid
		oneNotif.SourceId = int(group.Creator)
		oneNotif.Accepted = false
		oneNotif.CreatedAt = "not now"
		oneNotif.GroupId = int(oneinvite.GroupID)
		ResultArr = append(ResultArr, oneNotif)
	}
	followers, err := query.GetFollowers(context.Background(), int64(userid))
	if err != nil {
		log.Fatal(err)
	}
	for _, onefollower := range followers {
		if onefollower.Status == 0 {

			var oneNotif NotifStruct
			oneNotif.Label = "noti"
			oneNotif.Id = 0
			oneNotif.Type = "follow-req"
			oneNotif.TargetId = userid
			oneNotif.SourceId = int(onefollower.SourceID)
			oneNotif.Accepted = false
			oneNotif.CreatedAt = "not now"
			oneNotif.GroupId = 0
			ResultArr = append(ResultArr, oneNotif)
		}
	}
	fmt.Println("Most Important Thing: ", ResultArr)
	return ResultArr
}

func PrivateChatItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/private-chat-item"); err != nil {
			return
		}

		if r.Method == http.MethodGet {
			fmt.Printf("GET PrivateChatItemHandler\n")
			// Checks to find a user id in the url
			sourceId := r.URL.Query().Get("id")
			id, err := strconv.Atoi(sourceId)
			if err != nil {
				fmt.Println("Unable to convert to int")
			}

			foundId := false

			if sourceId != "" {
				foundId = true
			}
			fmt.Println("urlPrivateChatItem: ", r.URL)
			fmt.Printf("GET PrivateChatItemHandler %s\n", sourceId)

			// Declares the payload struct
			var Resp PrivateChatItemPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			if foundId {
				// ### GET CHAT ITEMS WITH SELF ID AS TARGET ID ###
				privateChatItems, err := query.GetPrivateChatItem(context.Background(), int64(id))

				if err != nil {
					fmt.Println(err)
					fmt.Println("Unable to find chat item")
				}

				for _, item := range privateChatItems {
					fmt.Println("chat item: ", item)
					// form the resp

					if err != nil {
						fmt.Println("Unable to find item")
					}

					var oneItem PrivateChatItemStruct

					oneItem.Id = int(item.ID)
					oneItem.SourceId = int(item.SourceID)
					oneItem.TargetId = int(item.TargetID)
					oneItem.ChatNoti = int(item.ChatNoti)
					oneItem.LastMsgAt = item.LastMsgAt.String()

					Resp.Data = append(Resp.Data, oneItem)
				}
			}

			// Marshals the response struct to a json object
			fmt.Println("p chat endpt Resp: ", Resp)
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
		} else {
			// Prevents all request types other than GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GroupChatItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-chat-item"); err != nil {
			return
		}

		userId := r.URL.Query().Get("userid")
		if userId == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		uId, err := strconv.Atoi(userId)

		if err != nil {
			fmt.Println("Unable to convert user ID")
		}

		switch r.Method {
		case http.MethodGet:
			// Declares the payload struct
			var Resp GroupChatItemPayload

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// ### get all group chat items ###

			chatItems, err := query.GetGroupChatNoti(context.Background(), int64(uId))

			if err != nil {
				fmt.Println("Unable to get group chat items")
			}

			for _, chatItem := range chatItems {
				var newChatItem GroupChatItemStruct

				newChatItem.Id = int(chatItem.ID)
				newChatItem.GroupId = int(chatItem.GroupID)
				newChatItem.UserId = int(chatItem.UserID)
				newChatItem.ChatNoti = int(chatItem.ChatNoti)
				newChatItem.LastMsgAt = chatItem.LastMsgAt.String()

				Resp.Data = append(Resp.Data, newChatItem)
			}

			// Marshals the response struct to a json object
			jsonResp, err := json.Marshal(Resp)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			// Sets the http headers and writes the response to the browser
			WriteHttpHeader(jsonResp, w)
			// case http.MethodPost:
			// 	// Declares the variables to store the group message details and handler response
			// 	var groupItem GroupChatItemStruct

			// 	// Decodes the json object to the struct, changing the response to false if it fails
			// 	err := json.NewDecoder(r.Body).Decode(&groupItem)
			// 	if err != nil {
			// 		fmt.Println(err.Error())
			// 	}

			// 	// ### CONNECT TO DATABASE ###

			// 	db := db.DbConnect()

			// 	query := crud.New(db)

			// 	// Check groupChatItem table to see if record exists

			// 	chatItem, err := query.GetGroupChatNotiByGroupId(context.Background(), int64(groupItem.GroupId))

			// 	if chatItem.GroupID != int64(groupItem.GroupId) {
			// 		// update chatItem last_msg_at column

			// 		_, err = query.UpdateGroupChatItem(context.Background(), crud.UpdateGroupChatItemParams{
			// 			LastMsgAt: time.Now(),
			// 			GroupID:   int64(groupItem.GroupId),
			// 		})

			// 		if err != nil {
			// 			fmt.Println("Unable to update group chat item")
			// 			return
			// 		}

			// 	} else {

			// 		// create group chat item

			// 		_, err = query.CreateGroupChatItem(context.Background(), crud.CreateGroupChatItemParams{
			// 			GroupID:   int64(groupItem.GroupId),
			// 			LastMsgAt: time.Now(),
			// 		})

			// 		if err != nil {
			// 			fmt.Println("Unable to create group chat item")
			// 			return
			// 		}
			// 	}

			// 	// set chatnoti to 0 (unseen) when chatitem is created or updated

			// 	_, err = query.UpdateGroupMemberChatNotiUnseen(context.Background(), int64(groupItem.GroupId))

			// 	if err != nil {
			// 		fmt.Println("Unable to update chatnoti for group memebers")
			// 		return
			// 	}

			// default:
			// 	// Prevents all request types other than POST and GET
			// 	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			// 	return
		}
	}
}

func GroupChatSeenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// EnableCors(&w)
		// Prevents the endpoint being called from other url paths
		if err := UrlPathMatcher(w, r, "/group-chat-seen"); err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:

		case http.MethodPost:
			// Declares the variables to store the group message details and handler response
			var groupMember GroupMemberStruct

			// Decodes the json object to the struct, changing the response to false if it fails
			err := json.NewDecoder(r.Body).Decode(&groupMember)
			if err != nil {
				fmt.Println(err.Error())
			}

			// ### CONNECT TO DATABASE ###

			db := db.DbConnect()

			query := crud.New(db)

			// set chatnoti to 1 (seen)

			_, err = query.UpdateGroupMemberChatNotiSeen(context.Background(), crud.UpdateGroupMemberChatNotiSeenParams{
				GroupID: int64(groupMember.GroupId),
				UserID:  int64(groupMember.UserId),
			})

			if err != nil {
				fmt.Println("Unable to update chatnoti for group members")
				return
			}

		default:
			// Prevents all request types other than POST and GET
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
