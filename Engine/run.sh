docker rm -f engine
docker run \
  -it \
  --name engine \
  -v "$(pwd)":/app:ro \
  engine:latest \
  /bin/bash -c "cmake -DCMAKE_BUILD_TYPE=Debug app; make; ./engine"