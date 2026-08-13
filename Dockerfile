# syntax=docker/dockerfile:1

# ---- Stage 1: build Go binary ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/compiler .

# ---- Stage 2: final runtime image ----
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    texlive-latex-base \
    perl \
    wget \
    ca-certificates \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# Debian's texlive packages are dpkg-managed, so tlmgr refuses to touch
# the system tree and requires an explicit user-mode tree (installs into
# $HOME/texmf, i.e. /root/texmf here, which is already on the default
# kpsewhich search path).
#
# Pin to the TeX Live 2025 historic archive, matching the release apt
# installed via texlive-latex-base on trixie. The live/rolling CTAN
# mirror tracks whatever the current year is and will refuse
# cross-release installs, and bookworm's texlive-latex-base is 2022,
# which also fails against this pinned archive.
RUN tlmgr init-usertree \
    && tlmgr option repository https://ftp.math.utah.edu/pub/tex/historic/systems/texlive/2025/tlnet-final \
    && tlmgr install \
        titlesec \
        hyperref \
        enumitem \
        fancyhdr \
        preprint \
        marvosym \
        fontawesome \
    && rm -rf /root/texmf/web2c/*.log

WORKDIR /app
COPY --from=builder /app/compiler /app/compiler

EXPOSE 8080

CMD ["/app/compiler"]