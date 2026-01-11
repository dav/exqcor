class Section < ApplicationRecord
  belongs_to :script, inverse_of: :sections
  has_many :props, inverse_of: :section, dependent: :delete_all
  has_many :sub_sections, -> { order(:ordering) }, inverse_of: :section, dependent: :delete_all
  has_many :character_sections, dependent: :delete_all
  has_many :characters, through: :character_sections

  after_create :build_first_sub_section, :add_vosd

  def next_ordering_index
    return 0 if sub_sections.empty?

    max = sub_sections.map(&:ordering).compact.max
    max ? max + 1 : 0
  end

  def ensure_vosd!
    add_vosd unless characters.include?(script.ensure_vosd!)
  end

  def add_vosd
    vosd = script.ensure_vosd!
    characters << vosd unless characters.include?(vosd)
  end

  def build_first_sub_section
    ss = SubSection.new(ordering: next_ordering_index)
    ss.section = self
    ss.save
  end

  def sorted_characters
    [script.VOSD] + characters.reject { |c| c == script.VOSD }
  end

  def writing_duration
    return 0 if sub_sections.empty?

    sub_sections.sum { |ss| ss.writing_duration || 0 }
  end

  def sort_name
    return "A" if name == "Opening Monologue"
    return "ZZZZZZZ" if name == "Closing Monologue"

    name
  end
end
