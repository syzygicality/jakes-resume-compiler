# syntax=docker/dockerfile:1

# ---- Stage 1: build Go binary ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/compiler .

# ---- Stage 2: TeX Live tree, shared by the format builder and the runtime ----
FROM debian:trixie-slim AS texlive

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

# ---- Stage 3: precompile the resume preamble into a LaTeX format file ----
#
# Dumping preamble.tex here lets per-request runs skip parsing hyperref, babel,
# fontawesome and friends. Kept in its own stage so the -ini run's .log/.aux
# droppings stay out of the runtime image.
FROM texlive AS fmt-builder

WORKDIR /fmt
COPY latex/preamble.tex ./preamble.tex

# nonstopmode matches compileHandler, letting the dump survive the same
# recoverable errors a live compile tolerates (hyperref 7.01p probing
# \IfDocumentMetadataT, undefined in Debian's 2024-11-01 kernel). That exits
# nonzero even on a good dump, so gate on preamble.fmt instead.
RUN pdftex -ini -interaction=nonstopmode -jobname=preamble "&pdflatex" preamble.tex; \
    test -s preamble.fmt || { cat preamble.log; echo "preamble.fmt was not dumped"; exit 1; }

# ---- Stage 4: final runtime image ----
FROM texlive AS runtime

# Give the format its own search path so pdflatex resolves it as plain
# `-fmt=preamble`; the trailing colon appends kpathsea's compiled-in default.
# Deliberately outside the texmf trees: those carry a `!!` prefix (ls-R only) in
# texmf.cnf, and kpathsea dedupes a plain entry against its `!!` twin, so a
# format dropped in one is invisible unless mktexlsr is also run.
ENV TEXFORMATS=/opt/texfmt:
COPY --from=fmt-builder /fmt/preamble.fmt /opt/texfmt/preamble.fmt

WORKDIR /app
COPY --from=builder /app/compiler /app/compiler

EXPOSE 8080

CMD ["/app/compiler"]