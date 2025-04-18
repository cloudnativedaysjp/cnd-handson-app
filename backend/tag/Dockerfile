FROM ruby:3.4.2-slim as base

WORKDIR /app
RUN apt-get -y update && apt-get install -y \
    build-essential\
    protobuf-compiler

COPY Gemfile Gemfile.lock ./

RUN bundle install

COPY . .

EXPOSE 4567
