PROJECT_NAME := "blacklist-example"

SUBDIRS := $(wildcard */.)

all: $(SUBDIRS) build

$(SUBDIRS):
	$(MAKE) -C $@

build:
	@go build -a -installsuffix cgo -o $(PROJECT_NAME) .

clean:
	@rm -f $(PROJECT_NAME)

.PHONY: all $(SUBDIRS)
