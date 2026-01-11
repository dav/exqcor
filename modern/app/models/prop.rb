class Prop < ApplicationRecord
  belongs_to :section, inverse_of: :props
end
