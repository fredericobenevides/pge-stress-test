build:
	docker build -t pge-stress-test --build-arg APP_PATH=./cmd/stresstest/main.go .

run-google:
	@echo "Running stress-test against google"
	docker run pge-stress-test --url=http://www.google.com --requests=10 --concurrency=2

run-amazon:
	@echo "Running stress-test against amazon"
	docker run pge-stress-test --url=http://www.amazon.com --requests=10 --concurrency=2

run: run-google