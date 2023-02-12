# forum 📚📚📚
[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger) 
Authors: Juss Märtson, Andrei Redi, Rain-Ander Laagus, Aleksandr Lavronenko, Joel Meeks

This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

This project is handled with Docker 
# Usage
#### Docker script 🐋
- ##### Try running bash script `bash start_Docker.sh`
```
~/j/go/forum > bash start_Docker.sh
Building docker image and container

🎉🎉🎉To stop server, press CTRL+C🎉🎉🎉
🎉🎉🎉   Press Enter to Continue  🎉🎉🎉
...
Starting server at port 8080
```
Docker image and container are deleted afterwards. 

#### By hand ⚒️
- ##### Try running command `go build`

    Now you should have file named `forum` 
```bash
~/j/go/forum > go build
~/j/go/forum > ls
Dockerfile   forum   main.go     sqlfunctions.go  web
...
```
- ##### Now, that you’ve created your executable you can run it `./forum`
```
~/j/go/forum > ./forum
Starting server at port 8080
~/j/go/forum >
```