create-local-queues:
	aws --endpoint-url=http://localhost:4566 sqs create-queue --region ap-southeast-2 --queue-name webhook-queue

purge-local-queues:
	aws --endpoint-url=http://localhost:4566 sqs purge-queue --region ap-southeast-2 --queue-url http://localhost:4566/000000000000/webhook-queue

local-env-up:
	docker-compose -f docker-compose.local.yml up -d

local-env-down:
	docker-compose -f docker-compose.local.yml down

test-unit:
	go test -short -race -v ./...

test:
	mkdir -p .coverage && \
	go test -race -coverprofile=.coverage/coverage.out -v ./...

coverage:
	go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html && \
	open .coverage/coverage.html

clean:
	rm -rf .coverage