FROM alpine:3.20.3

LABEL org.opencontainers.image.source=https://github.com/freepikcompany/gommitizen
LABEL org.opencontainers.image.description="A commitizen implementation for Go with multi-project support"
LABEL org.opencontainers.image.licenses=MIT

RUN apk --no-cache add \
    ca-certificates \
    openssh-client \
    git \
    bash

RUN addgroup -S gommitizen && adduser -S gommitizen -G gommitizen

COPY gommitizen /usr/local/bin/gommitizen

# Mount the source code here
WORKDIR /source

USER gommitizen

RUN echo "git config --global --add safe.directory /source" >> /home/gommitizen/.bashrc && \
    echo 'git config --global user.name "$GIT_USER_NAME"' >> /home/gommitizen/.bashrc && \
    echo 'git config --global user.email "$GIT_USER_EMAIL"' >> /home/gommitizen/.bashrc

# Set the entrypoint to start the SSH agent
ENTRYPOINT ["/bin/bash", "-c", "source /home/gommitizen/.bashrc && /usr/local/bin/gommitizen $@", "--"]
CMD [""]
