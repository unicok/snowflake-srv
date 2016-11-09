FROM alpine:3.2
ADD snowflake /snowflake
ENTRYPOINT [ "/snowflake" ]