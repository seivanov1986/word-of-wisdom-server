INSTALL:

cd docker
cp env.sample env
docker-compose build
docker-compose up -d
