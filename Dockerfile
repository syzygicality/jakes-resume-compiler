# ---- Stage 1: build Python deps ----
FROM python:3.13-slim AS builder

COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

WORKDIR /app

COPY pyproject.toml uv.lock main.py ./

RUN uv sync --frozen

# ---- Stage 2: final runtime image ----
FROM python:3.13-slim

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
# kpsewhich search path). Installing directly in the final stage (rather
# than a separate stage + multi-path COPY) avoids having to reconstruct
# Debian's scattered texlive layout (binaries, symlinks, and config
# spread across /usr/bin, /usr/share/texlive, /etc/texmf, /var/lib/texmf)
# by hand.
#
# Pin to the TeX Live 2025 historic archive, matching the release apt
# installed via texlive-latex-base. The live/rolling CTAN mirror tracks
# whatever the current year is and will refuse cross-release installs.
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

COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv
WORKDIR /app
COPY --from=builder /app /app

CMD ["uv", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
