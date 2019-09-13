FROM alpine

# Define working directory.
WORKDIR /opt/
ADD app /opt/

# Copy configs
ADD resources /opt/resources

EXPOSE 50051 8081

# Define default command.
CMD ["./app", "-env=docker"]
