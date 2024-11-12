# Use the base image
FROM mcr.microsoft.com/devcontainers/universal:2

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg2 \
    lsb-release \
    software-properties-common

# Add PostgreSQL APT repository
RUN echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -

# Install PostgreSQL 13 client tools
RUN apt-get update && apt-get install -y postgresql-client-13 && rm -rf /var/lib/apt/lists/*