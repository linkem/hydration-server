#only host build
#docker build -t hydration-go-backend:latest --name hydration-go-backend


#docker buildx create --name mybuilder
#docker buildx use mybuilder
#https://docs.docker.com/docker-for-mac/multi-arch/
#multi platform build
#docker buildx build --platform linux/amd64,linux/arm64 -t mlinke/hydration-go-backend:latest --push .