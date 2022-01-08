docker rm -f engine
docker run \
  -it \
  --name engine \
  -v "$(pwd)":/app:ro \
  engine:latest \
  /bin/bash -c "cmake app; make; ./engine"