class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def create_linked_list(values):
    if not values:
        return None

    head = ListNode(values[0])
    current = head

    for value in values[1:]:
        current.next = ListNode(value)
        current = current.next

    return head

def print_linked_list(head):
    current = head
    while current:
        print(current.val, end=" -> ")
        current = current.next

    print("None")

def reverse_linked_list(head):
    prev = None
    curr = head

    while curr:
        next_node = curr.next   # save the next node
        curr.next = prev        # reverse the current pointer
        prev = curr             # move prev node to current node
        curr = next_node        # move current to the next node
    return prev     # this is the head now

# Setup
values = [1,2,3,4,5]
head = create_linked_list(values)
print_linked_list(head)

reverse_head = reverse_linked_list(head)
print_linked_list(reverse_head)
