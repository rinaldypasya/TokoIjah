SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

clean:
	@rm -f $(TARGET)

build:
	@go build -o $(TARGET)

rebuild: clean build