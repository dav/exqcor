class Character < ActiveRecord::Base
  belongs_to :play, :inverse_of => :characters
  attr_accessible :description, :name
  has_many :lines, :inverse_of => :sub_section
end
