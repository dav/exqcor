class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play, :dependent => :delete_all
  has_many :sections, :inverse_of => :play, :dependent => :delete_all
end
