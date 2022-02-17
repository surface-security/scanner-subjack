build:
	docker build -t tmp-$(notdir $(CURDIR)) .
