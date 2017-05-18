FROM scratch

ENV GITHUBINT_SERVICE_PORT 8080
ENV GITHUBINT_BRANCH "workshop/"
ENV GITHUBINT_TOKEN "Webhook secret is in integration settings on Github"
ENV GITHUBINT_PRIV_KEY "Private key is in integration settings on Github"
ENV GITHUBINT_INTEGRATION_ID "Integration ID is in it's settings on Github"

ENV CICD_SERVICE_HOST "localhost"
ENV CICD_SERVICE_PORT "8080"
ENV USERMAN_SERVICE_HOST "localhost"
ENV USERMAN_SERVICE_PORT "8080"

EXPOSE $GITHUBINT_SERVICE_PORT

COPY github-integration /

CMD ["/github-integration"]