default: challenge

TODAY := $(shell date +%Y%m%d)

challenge:
	@read -p "Enter the challenge name: " challenge; \
		git switch main; \
		git pull; \
		git switch -c $(TODAY)-$$challenge; \
		mkdir -p $(TODAY)-$$challenge; \
		cd $(TODAY)-$$challenge && touch Readme.md main.py
