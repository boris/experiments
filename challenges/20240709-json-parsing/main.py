import urllib.request
import json

def get_object():
    try:
        r = urllib.request.urlopen("https://jsonplaceholder.typicode.com/todos/1aaaaa")
        data = json.loads(r.read())
        return data
    except urllib.error.URLError as e:
        print(f"Failed to fetch data: {e}")
        return None
    except json.JSONDecodeError:
        print("Failed to decode JSON")
        return None


def parse_data(data):
    if data:
        title = data.get('title', 'No title found')
        print(f"Title: {title}")
    else:
        print("No data to parse")



def main():
    data = get_object()
    parse_data(data)


if __name__ == "__main__":
    main()
