default: challenge

TODAY := $(shell date +%Y%m%d)

challenge:
	git switch -c $(TODAY)-challenge
	mkdir -p $(TODAY)
	cd $(TODAY) && touch Readme.md main.py
