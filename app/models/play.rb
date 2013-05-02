class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play
end
