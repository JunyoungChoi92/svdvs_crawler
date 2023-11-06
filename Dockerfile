# Use the LXDE VNC image as a base
FROM dorowu/ubuntu-desktop-lxde-vnc:bionic-lxqt

# Set environment variables for Go installation
# following by 2022 spo-vdvs-system. So it can be changed anytime.
ENV GOLANG_VERSION 1.21.2

# Install Go
RUN apt-get update && \
    apt-get install -y wget git && \
    wget https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm go$GOLANG_VERSION.linux-amd64.tar.gz

# Set environment variables for the Go project
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Create the directory for the Go project
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install Google Chrome
# Processor's architecture version can be changed anytime by it's running environment.
RUN apt-get update && apt-get install -y \
    dpkg \
    fonts-liberation \
    libu2f-udev \
    libvulkan1 \ 
    && wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && dpkg -i google-chrome-stable_current_amd64.deb

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
CMD ["./newcrawler"]


