
build-all:
	SERVICE=$(service) ./go-mod-tidy.sh
	SERVICE=$(service) ./gen_desc.sh
	SERVICE=$(service) ./gen-swagger-doc.sh
	SERVICE=$(service) ./build-all.sh
	SERVICE=$(service) ./gen-swagger-doc.sh


docker-build:
	SERVICE=$(service) ./docker-build.sh


project-init:
	git clone git@gitlab.com:t8322/boilerplate.git ../$(service)
	rm -rf ../$(service)/.git

pull-all:
	chmod +x ./git-pull.sh
	SERVICE=$(service) ./git-pull.sh

logs:
	chmod +x ./log-checker.sh
	service=$(service) ./log-checker.sh

