export DOCKER_BUILDKIT=1
sudo docker build --target backend-builder -f Dockerfile -t  backend-builder .
sudo docker build --target frontend-builder -f ./frontend/Dockerfile -t frontend-builder .
sudo docker build --target final -f Dockerfile -t myoauth .
