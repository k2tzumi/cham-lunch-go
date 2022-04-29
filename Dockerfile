# syntax=docker/dockerfile:1
FROM golang:1.18.1-bullseye

RUN apt-get update
RUN apt-get install -y --no-install-recommends \
      build-essential \
      git \
      gcc \
      zsh \
      less \
      vim \
      wget
ENV LANG=ja_JP.UTF-8 \
    LANGUAGE=ja_JP:ja\
    LC_ALL=ja_JP.UTF-8
RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
RUN wget -O ~/vsls-reqs https://aka.ms/vsls-linux-prereq-script && \
    chmod +x ~/vsls-reqs && \
    ~/vsls-reqs
RUN groupadd -g 999 devuser && \
    useradd -m -u 999 -g devuser -s /bin/zsh devuser

COPY ./entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

USER devuser
RUN go install golang.org/x/tools/gopls@latest
RUN go install mvdan.cc/gofumpt@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/ramya-rao-a/go-outline@latest
WORKDIR /home/devuser/workspace
ENTRYPOINT ["entrypoint.sh"]
CMD ["tail", "-f", "/dev/null"]
