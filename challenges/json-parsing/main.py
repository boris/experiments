import urllib.request
import json

def get_object():
    response = urllib.request.urlopen("https://jsonplaceholder.typicode.com/posts")
    return json.loads(response.read())


def parse_data(data):
    extracted_data = []
    user_post_count = {}

    for item in data:
        id_ = item.get("id")
        title = item.get("title")
        user_id = item.get("userId")

        extracted_data.append({
            "id": id_,
            "title": title,
            "user_id": user_id
        })

        # Count posts by user
        if user_id in user_post_count:
            user_post_count[user_id] += 1
        else:
            user_post_count[user_id] = 1

    return extracted_data, user_post_count


def print_results(extracted_data, user_post_count):
    for data in extracted_data:
        print(f"ID: {data['id']}, Title: {data['title']}, User ID: {data['user_id']}")

    print("\nPost counts by User ID:")
    for user_id, count in user_post_count.items():
        print(f"User ID: {user_id}, Post Count: {count}")


def main():
    data = get_object()
    extracted_data, user_post_count = parse_data(data)
    print_results(extracted_data, user_post_count)



if __name__ == "__main__":
    main()
