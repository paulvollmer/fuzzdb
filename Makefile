build: go js

go:
	@go run generator.go -lang go

js:
	@go run generator.go -lang js

clean:
	rm -f attack/**/*.go
	rm -f attack/**/*.js
	rm -f attack/**/**/*.go
	rm -f attack/**/**/*.js
	rm -f discovery/**/*.go
	rm -f discovery/**/*.js
	rm -f discovery/**/**/*.go
	rm -f discovery/**/**/*.js
	rm -f regex/*.go
	rm -f regex/*.js
	rm -f wordlists-misc/*.go
	rm -f wordlists-misc/*.js
	rm -f wordlists-user-passwd/*.go
	rm -f wordlists-user-passwd/*.js
	rm -f wordlists-user-passwd/**/*.go
	rm -f wordlists-user-passwd/**/*.js
