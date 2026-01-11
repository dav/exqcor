class Line < ApplicationRecord
  belongs_to :sub_section, inverse_of: :lines
  belongs_to :character, inverse_of: :lines
end
