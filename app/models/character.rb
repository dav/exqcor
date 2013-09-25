class Character < ActiveRecord::Base
  before_destroy :confirm_not_vosd

  belongs_to :play, :inverse_of => :characters
  attr_accessible :description, :name
  has_many :lines, :inverse_of => :sub_section
  has_and_belongs_to_many :sections
  validates_uniqueness_of :name, scope: :play_id
  
  def confirm_not_vosd
    return false if self.name == 'VOSD'
    true
  end
end
