Requirements:

1. Use curl to fetch JSON data from a provided URL.
1. Parse the JSON data to extract the value of a specific key.
1. Print the extracted value.

URL for JSON Data:
You can use this placeholder URL for the challenge: `https://jsonplaceholder.typicode.com/todos/1`

Steps to Complete:

1. Fetch the JSON data using curl.
1. Parse the JSON data using a tool of your choice (e.g., jq for shell scripting, or JSON libraries in Python, JavaScript, etc.).
1. Extract the value of the key title from the JSON data.
1. Print the extracted value.

Example JSON Response:
```json
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}

```

Output:
```
Title: delectus aut autem
```
