docker rmi nandanurseptama/golang-grafana-otel/user-service:latest
docker rmi nandanurseptama/golang-grafana-otel/auth-service:latest 
docker rmi nandanurseptama/golang-grafana-otel/frontend-service:latest 
docker build -t nandanurseptama/golang-grafana-otel/user-service:latest services/user/
docker build -t nandanurseptama/golang-grafana-otel/auth-service:latest services/auth/
docker build -t nandanurseptama/golang-grafana-otel/frontend-service:latest services/frontend/