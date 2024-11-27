FROM gcc:latest

WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y \
    cmake \
    libboost-all-dev \
    libssl-dev \
    libcurl4-openssl-dev \
    nlohmann-json3-dev

COPY . .

# Build
RUN mkdir build && cd build && \
    cmake .. && \
    make -j$(nproc)

EXPOSE 9000

CMD ["./build/auth_service"]
