# forum 📚📚📚
[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger) 
Authors: Juss Märtson, Andrei Redi, Rain-Ander Laagus, Aleksandr Lavronenko, Joel Meeks

This bonus task consists of making our forum more interactive :

- Comment and Post Removal/Editing.
- Saves all user activity and displays it on "Account" page.
- Displays notifications from other users, who liked/disliked/commented your post or comment.

## Useful links
### audit page: https://github.com/01-edu/public/blob/master/subjects/forum/advanced-features/audit.md

# Usage
#### The Easy Way
```
~/j/go/forum-forum-advanced-features > go run .
```

#### Docker script 🐋
- ##### Try running bash script `bash start_Docker.sh`
```
~/j/go/forum-advanced-features > bash start_Docker.sh
Building docker image and container

🎉🎉🎉To stop server, press CTRL+C🎉🎉🎉
🎉🎉🎉   Press Enter to Continue  🎉🎉🎉
...
Starting server at port 8080
```
Docker image and container are deleted afterwards. 

# Good to know
- All your account activity (Posts, comments, likes, dislikes) are under "Account" page.
- All notifications from other users are under "Activity" page, you can go there by clicking the bell button on top right corner. 
- Notification will disappear once you go on the post it's from.

# Accounts for testing
- User Name: isabella
- Password: Isabella200