list:
	@echo "Available targets:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

build-docker-image:
	docker build -f ./scripts/Dockerfile -t foxfox-codelabs-preview:latest .
	docker tag foxfox-codelabs-preview:latest gcr.io/foxfox-learn/foxfox-codelabs-preview:staging
	docker push gcr.io/foxfox-learn/foxfox-codelabs-preview:staging

deploy-cloud-run:
	gcloud run deploy foxfox-codelabs-preview --image gcr.io/foxfox-learn/foxfox-codelabs-preview:staging --platform managed --allow-unauthenticated --port=4000 --region=asia-northeast1 --update-env-vars FOXFOX_PLATFORM=gcs,FOXFOX_CONFIG_BUCKET=foxfox-gcs,FOXFOX_CONFIG_PATH=.credential/config.yml
