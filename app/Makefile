SRC := $(shell find ./ -name "*.go")
TARGET := app
OBJ :=

.PHONY: build
build: $(TARGET) $(OBJ)

$(TARGET): $(SRC)
	go build

.PHONY: clean
clean:
	rm $(TARGET) $(OBJ)
