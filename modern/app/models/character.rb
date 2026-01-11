class Character < ApplicationRecord
  before_destroy :confirm_not_vosd

  belongs_to :script, inverse_of: :characters
  belongs_to :actor, inverse_of: :characters, optional: true

  has_many :lines, inverse_of: :character, dependent: :delete_all
  has_many :character_sections, dependent: :delete_all
  has_many :sections, through: :character_sections

  enum :role, { character: "character", vosd: "vosd" }

  validates :name, uniqueness: { scope: :script_id }
  validates :role, inclusion: { in: roles.keys }

  def confirm_not_vosd
    throw(:abort) if vosd?

    true
  end
end
