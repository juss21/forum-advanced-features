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
docker container run -p 8080:8080 --detach --name dockerize-container dockerize-image

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




