class Writer < ApplicationRecord
  has_many :sub_sections, inverse_of: :writer, dependent: :nullify
end
