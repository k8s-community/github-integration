FROM alpine:3.5

ENV GITHUBINT_SERVICE_PORT 8080
ENV GITHUBINT_BRANCH "release/"

ENV GITHUBINT_TOKEN "Webhook secret is in integration settings on Github"
ENV GITHUBINT_PRIV_KEY "Private key is in integration settings on Github"
ENV GITHUBINT_INTEGRATION_ID "Integration ID is in it's settings on Github"

# DB parameters
ENV COCKROACHDB_PUBLIC_SERVICE_HOST localhost
ENV COCKROACHDB_PUBLIC_SERVICE_PORT 26257
ENV COCKROACHDB_USER githubint
ENV COCKROACHDB_PASSWORD githubint
ENV COCKROACHDB_NAME github_integration

ENV CICD_BASE_URL http://k8s-build-01:8080
ENV USERMAN_BASE_URL https://services.k8s.community/user-manager

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY github-integration /

CMD ["/github-integration"]

EXPOSE $GITHUBINT_SERVICE_PORT
