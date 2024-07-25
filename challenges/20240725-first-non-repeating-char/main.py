def check_first_non_rep(string):
    # string.count will return the number of occurrences of each string.
    # example:
    # >>> "aaabb".count("a")
    # 3
    #
    # the code will iterate over the string counting each character. if the
    # count is == 1 means that the character is not repeated
    for i in range(len(string)):
        if string.count(string[i]) == 1:
            return string[i]
    return "_"



if __name__ == '__main__':
    print(check_first_non_rep("swiss"))
    print(check_first_non_rep("releveler"))
    print(check_first_non_rep("aabbcc"))
