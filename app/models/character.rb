class Character < ActiveRecord::Base
  belongs_to :play, :inverse_of => :characters
  attr_accessible :description, :name
end
