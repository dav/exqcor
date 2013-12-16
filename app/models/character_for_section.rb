class CharacterForSection < ActiveRecord::Base
  set_table_name 'characters_sections' 
  belongs_to :character
  belongs_to :section
end