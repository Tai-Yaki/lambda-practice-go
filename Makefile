STACK_NAME := lambda-practice-go
TEMPLATE_FILE := template.yaml
SAM_FILE := sam.yaml

build: build-login build-user
.PHONY: build

build-login:
	GOARCH=amd64 GOOS=linux go build -o artifact/login ./handler/login
.PHONY: build-login

build-user: build-user-create build-user-show build-user-index build-user-update build-user-delete
.PHONY: build-user

build-user-create:
	GOARCH=amd64 GOOS=linux go build -o artifact/user/create ./handler/user/create
.PHONY: build-user-create

build-user-show:
	GOARCH=amd64 GOOS=linux go build -o artifact/user/show ./handler/user/show
.PHONY: build-user-show

build-user-index:
	GOARCH=amd64 GOOS=linux go build -o artifact/user/index ./handler/user/index
.PHONY: build-user-index

build-user-update:
	GOARCH=amd64 GOOS=linux go build -o artifact/user/update ./handler/user/update
.PHONY: build-user-update

build-user-delete:
	GOARCH=amd64 GOOS=linux go build -o artifact/user/delete ./handler/user/delete
.PHONY: build-user-delete

deploy: build
	sam package --template-file $(TEMPLATE_FILE) --s3-bucket $(STACK_NAME) --output-template-file $(SAM_FILE)
	sam deploy --template-file $(SAM_FILE) --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --parameter-override LinkTableName=$(LINK_TABLE)
	echo API endpoint URL for Prod environment:
		aws cloudformation desctibe-stacks --stack-name $(STACK_NAME) --query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' --output text
.PHONY: deploy

delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
	aws s3 rm "s3://$(STACK_BUCKET)" --recursive
	aws s3 rb "s3://$(STACK_BUCKET)"
.PHONY: delete

test:
	go test -v ./...
