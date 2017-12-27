
APPS=dns validate cleanup

app:
	$(foreach exe,$(APPS),\
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/letsencrypt-$(exe)/letsencrypt-$(exe) cmd/letsencrypt-$(exe)/main.go ; \
	)
install:
	$(foreach exe,$(APPS),\
		install -m 755 cmd/letsencrypt-$(exe)/letsencrypt-$(exe) /usr/local/bin/; \
	)

uninstall:
	$(foreach exe,$(APPS),\
		rm -f  /usr/local/bin/letsencrypt-$(exe); \
	)

clean:
	$(foreach exe,$(APPS),rm -f cmd/letsencrypt-$(exe)/letsencrypt-$(exe);)

image:
	docker build --rm -t helber/letsencrypt-dns .
run:
	docker run -it helber/letsencrypt-dns
