FROM python:3.12.3

RUN mkdir /sandbox \
    && useradd --uid 3456 -M -d /sandbox sandbox \
    && chown -R sandbox:sandbox /sandbox