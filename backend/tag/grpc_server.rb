# frozen_string_literal: true

require_relative './lib/tag_services_pb'

class ServerImpl < TagService::Service
  def getTag
    tag = Tag.new(
        id:,
        name:,
        text_color:,
        background_color:,
    )

    TagResponse.new(error: nil, tag:)
  end
end