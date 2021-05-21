.PHONY: build start-dynamodb create-product-table run-dynamodb stop-dynamodb list-tables query-products

build:
	sam build

start-dynamodb: stop-dynamodb run-dynamodb create-product-table

create-product-table:
	sleep 5 && aws dynamodb create-table --cli-input-json file://config/create-product-table.json --endpoint-url http://localhost:8000

run-dynamodb:
	docker run -d -p 8000:8000 --name dynamodb amazon/dynamodb-local

stop-dynamodb:
	docker stop dynamodb || true && docker rm dynamodb || true

list-tables:
	aws dynamodb list-tables --endpoint-url http://localhost:8000

query-products:
	aws dynamodb query --table-name Product --endpoint-url http://localhost:8000