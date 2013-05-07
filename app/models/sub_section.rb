class SubSection < ActiveRecord::Base
  attr_accessible :ordering, :play_id
  belongs_to :section, :inverse_of => :sub_sections
  has_many :lines, :inverse_of => :sub_section, :dependent => :delete_all
end
