mock: install-mock
	find . -type f -name *_mock.go -delete
	mockgen -source=./pkg/services/auth/service.go -destination=./pkg/services/auth/service_mock.go -package=auth

install-mock:
	@REQUIRED_VERSION=$$(grep go.uber.org/mock go.mod | sed 's/.* //'); \
	if which mockgen > /dev/null 2>&1; then \
		CURRENT_VERSION=$$(mockgen -version); \
		if [ "$$REQUIRED_VERSION" != "$$CURRENT_VERSION" ]; then \
			echo Version Mismatched, install mockgen executable...; \
			go install go.uber.org/mock/mockgen@$$REQUIRED_VERSION || exit 1; \
			echo Install mockgen successful; \
		fi; \
	else \
		echo Install mockgen executable...; \
		go install go.uber.org/mock/mockgen@$$REQUIRED_VERSION || exit 1; \
		echo Install mockgen successful; \
	fi;
