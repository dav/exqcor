class Actor < ApplicationRecord
  has_many :characters, inverse_of: :actor, dependent: :delete_all
end
