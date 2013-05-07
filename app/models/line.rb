class Line < ActiveRecord::Base
  attr_accessible :character_id, :ordering, :sub_section_id, :text
  belongs_to :sub_section, :inverse_of => :lines
  belongs_to :character, :inverse_of => :lines
end
