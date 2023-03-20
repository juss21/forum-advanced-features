 #!/bin/bash

echo "Building docker image and container"
echo 
echo "🎉🎉🎉To stop server, press CTRL+C🎉🎉🎉"
echo "🎉🎉🎉   Press Enter to Continue  🎉🎉🎉"

echo ""
echo ""
echo ""

read "part"




docker image build -f Dockerfile -t dockerize-image .
docker container run -p 8080:8080 --detach  --name dockerize-container dockerize-image

echo ""
echo "Server is running in http://www.localhost:8080/"
echo ""
echo "Here on can acces to files in container. "
echo "In case you need check data from database"
echo "database filename is database.db"
echo "To exit, write \"exit\""
echo ""
echo ""

docker exec -it dockerize-container /bin/bash



echo ""
echo ""
echo ""

echo "Removing Docker image and container from your system"
echo ""
echo ""
docker rm -f dockerize-container
docker image rm dockerize-image


echo "🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉"
echo "DONE"




