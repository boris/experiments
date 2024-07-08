Objective: Write a script that fetches data from a given URL, parses the JSON
response, and extracts specific information. Then, perform a simple manipulation
on the extracted data.

Task:
URL to Use: https://jsonplaceholder.typicode.com/posts
- Fetch data from the given URL.
- Extract the following fields: id, title, and userId.
- Count the number of posts by each user (userId).
- Output the extracted data and the count of posts by each user.

Example Output:
```
ID: 1, Title: "Title1", User ID: 1
ID: 2, Title: "Title2", User ID: 1
ID: 3, Title: "Title3", User ID: 2
...
User ID 1 has 10 posts
User ID 2 has 7 posts
...
```

Requirements:
- The script should be written in Python.
- Use appropriate libraries for HTTP requests and JSON parsing.
- Ensure the script handles potential errors (e.g., network issues, invalid JSON).
