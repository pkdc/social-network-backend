# Social Network Application

This is the backend of a Facebook-like social network application built using React for the frontend and Golang for the backend.

It is hosting on https://socialnetwork123.netlify.app/

## Features

The application includes the following features:

- User registration and login with sessions and cookies for authentication
- User profiles with public and private settings
- Ability to follow/unfollow other users
- Creation of posts and comments with privacy settings
- Creation of groups with invitations and requests
- Real-time private and group messaging using Websocket
- Notifications for following requests, group invitations and requests, and group events using Websocket
- Dockerized deployment with separate backend and frontend images

## Technologies Used

- Frontend: React (JavaScript framework)
- Backend: Golang (Go programming language)
- Database: SQLite
- WebSocket: Gorilla WebSocket package
- Migration: golang-migrate package
- Authentication: Sessions and cookies
- Containerization: Docker

## Chat rules

- The user should be able to send private messages to users that he/she is following
- The user that the message was sent to, will receive the message instantly, if he/she is following the user that sent the message (shown under "Users You Are Following:") or if the user has a public profile (shown under "Other Users:")
- if the user is not following the user that sent the message and if the user doesn't have a public profile, the message is stored the database, and it will be shown when he/she starts following the user that sent the message or if the user switches his/her profile to public 
- If a public user switch his/her profile back to private, he/she can still continue on the existing conversations (even if he/she didn't send any messages)
- Groups should have a common chat room, so if a user is a member of the group he/she should be able to send and receive messages to this group chat.
