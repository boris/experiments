.PHONY: freeze

freeze:
	pip freeze > requirements.txt

run:
	uvicorn main:app --reload

json_test:
	curl -X POST -H "Content-Type: application/json" -d @test.json http://localhost:8000
