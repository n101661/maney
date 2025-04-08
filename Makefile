mock: install-mock
	find . -type f -name *_mock.go -delete
	mockgen -source=./server/services/auth/service.go -destination=./server/services/auth/service_mock.go -package=auth

models: install-openapi-codegen
	@find . -type f -name *_gen.go -delete; \
	mkdir -p ./server/models; \
	oapi-codegen --generate=types --package=models -o=./server/models/models_gen.go openapi.yaml; \
	echo Done;

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

install-openapi-codegen:
	@REQUIRED_VERSION=$$(grep github.com/oapi-codegen/oapi-codegen/v2 go.mod | sed 's/.* //'); \
	if which oapi-codegen > /dev/null 2>&1; then \
		CURRENT_VERSION=$$(oapi-codegen -version | grep -E v[0-9]+\\.[0-9]+\\.[0-9]+); \
		if [ "$$REQUIRED_VERSION" != "$$CURRENT_VERSION" ]; then \
			echo Version Mismatched, install oapi-codegen executable...; \
			go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen || exit 1; \
			echo Install oapi-codegen successful; \
		fi; \
	else \
		echo Install oapi-codegen executable...; \
		go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen || exit 1; \
		echo Install oapi-codegen successful; \
	fi;
