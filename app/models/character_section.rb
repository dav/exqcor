class CharacterSection < ActiveRecord::Base
  attr_accessible :on_stage
  belongs_to :character, :inverse_of => :sections
  belongs_to :section, :inverse_of => :characters
end