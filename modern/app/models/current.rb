class Current < ActiveSupport::CurrentAttributes
  attribute :session, :user, :api_token
  delegate :user, to: :session, allow_nil: true
end
