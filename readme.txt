go mod tidy

clear; go build -o app . &&\
docker build -t go-redmine-ish:latest -f DockerfileShort . &&\
docker tag go-redmine-ish localhost:32000/go-redmine-ish-golang-app:latest &&\
docker push localhost:32000/go-redmine-ish-golang-app:latest &&\
microk8s kubectl rollout restart deploy go-redmine-ish-golang-app -n go-redmine-ish
