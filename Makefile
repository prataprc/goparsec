SUBDIRS := json expr

build:
	go build ./...
	@for dir in $(SUBDIRS); do \
		echo $$dir "..."; \
		$(MAKE) -C $$dir build; \
	done

test:
	go test -v -race -timeout 4000s -test.run=. -test.bench=. -test.benchmem=true ./...
	@for dir in $(SUBDIRS); do \
		echo $$dir "..."; \
		$(MAKE) -C $$dir test; \
	done

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm -rf coverage.out
	@for dir in $(SUBDIRS); do \
		echo $$dir "..."; \
		$(MAKE) -C $$dir test; \
	done

vet:
	go vet ./...

lint:
	golint ./...

.PHONY: build $(SUBDIRS)
