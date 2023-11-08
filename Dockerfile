# Use the LXDE VNC image as a base
FROM dorowu/ubuntu-desktop-lxde-vnc:focal

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 4EB27DB2A3B88B8B
RUN apt-get update -y

RUN apt-get install -y -q curl
RUN apt-get install -y -q wget
RUN apt-get install -y -q gnupg
RUN apt-get install -y -q ca-certificates
RUN apt-get install -y xvfb
RUN apt-get -y install xorg xvfb gtk2-engines-pixbuf
RUN apt-get -y install dbus-x11 xfonts-base xfonts-100dpi xfonts-75dpi xfonts-cyrillic xfonts-scalable
RUN apt-get -y install imagemagick x11-apps

RUN Xvfb -ac :99 -screen 0 1280x1024x16 & export DISPLAY=:99

# Timezone 셋팅
RUN ln -fs /usr/share/zoneinfo/Asia/Seoul /etc/localtime

# Set environment variables for Go installation
# following by 2022 spo-vdvs-system. So it can be changed anytime.
ENV GOLANG_VERSION 1.21.2

# Install Go
RUN apt-get update && \
    apt-get install -y wget git && \
    wget https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm go$GOLANG_VERSION.linux-amd64.tar.gz

# chrome 설치
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt-get install -y ./google-chrome-stable_current_amd64.deb

# pm2 설치
RUN apt-get install -y nodejs && \
    apt-get install -y npm && \
    npm i -g pm2

# Set environment variables for the Go project
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Create the directory for the Go project
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Clone the repository
RUN git clone https://github.com/JunyoungChoi92/svdvs_crawler.git $GOPATH/src/crawler

# Set the working directory to the crawler's cmd directory
WORKDIR $GOPATH/src/crawler/cmd

# Install dependencies and build the Go program
RUN go mod tidy
RUN go build -o newcrawler .

# Set the display environment variable for ChromeDP
ENV DISPLAY :1

# When the container starts, run the compiled Go program
CMD ["./crawler"]


