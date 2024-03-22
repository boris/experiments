.PHONY: freeze

freeze:
	pip freeze > requirements.txt

run:
	uvicorn main:app --reload

json_test:
	curl -X POST -H "Content-Type: application/json" -d @example_data.json http://localhost:8000

check:
	curl -X POST -H "Content-Type: application/json" -d @data.json http://API_URL
