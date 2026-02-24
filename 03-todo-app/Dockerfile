FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libpq5 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY .artifacts/libchameleon.so /usr/local/lib/
COPY .artifacts/chameleon.h /usr/local/include/

RUN ldconfig

# en el futuro:
# COPY chameleond /usr/local/bin/chameleond

CMD ["bash"]
