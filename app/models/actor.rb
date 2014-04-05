class Actor < ActiveRecord::Base
  attr_accessible :name

  has_many :characters, :inverse_of => :actor
end
