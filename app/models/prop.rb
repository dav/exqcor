class Prop < ActiveRecord::Base
  attr_accessible :description, :name, :play_id
  belongs_to :section, :inverse_of => :props
end
