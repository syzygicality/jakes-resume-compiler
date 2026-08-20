# ---- Stage 1: build Go binary ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/compiler .

# ---- Stage 2: pinned TeX Live base, shared by the format builder and the runtime ----
#
# Built from latex/Dockerfile.base and pushed to GHCR. Avoids hitting a TeX
# Live historic mirror on every build. Only rebuilt when the vendored
# package set changes. See latex/Dockerfile.base for the tlmgr install list.
FROM ghcr.io/syzygicality/jrc-texlive-base:2025 AS texlive

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