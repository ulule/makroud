test:
	@(echo "-> Running unit tests...")
	@(SQLXX_DISABLE_CACHE=1 go test -v . ./reflekt)
