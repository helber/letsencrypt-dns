
APPS=letsencrypt-dns letsencrypt-validate letsencrypt-cleanup checkcert oc-patch-route

app:
	$(foreach exe,$(APPS),\
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/$(exe)/$(exe) cmd/$(exe)/main.go ; \
	)
install:
	$(foreach exe,$(APPS),\
		install -m 755 cmd/$(exe)/$(exe) /usr/local/bin/; \
	)

uninstall:
	$(foreach exe,$(APPS),\
		rm -f  /usr/local/bin/$(exe); \
	)

clean:
	$(foreach exe,$(APPS),rm -f cmd/$(exe)/$(exe) cmd/$(exe)/debug ;)

image:
	docker build --rm -t helber/letsencrypt-dns .
run:
	docker run -it helber/letsencrypt-dns
