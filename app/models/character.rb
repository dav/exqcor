class Character < ActiveRecord::Base
  belongs_to :play, :inverse_of => :characters
  attr_accessible :description, :name
  has_many :lines, :inverse_of => :sub_section
  has_and_belongs_to_many :sections
  validates_uniqueness_of :name, scope: :play_id
end
