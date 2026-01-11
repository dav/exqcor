class SubSection < ApplicationRecord
  belongs_to :section, inverse_of: :sub_sections
  belongs_to :writer, inverse_of: :sub_sections, optional: true
  has_many :lines, -> { order(:ordering) }, inverse_of: :sub_section, dependent: :delete_all

  def next_ordering_index
    return 0 if lines.empty?

    max = lines.map(&:ordering).compact.max
    max ? max + 1 : 0
  end

  def next_section
    sss = section.sub_sections
    sss.each_with_index do |sub_section, i|
      if sub_section == self
        return sss[i + 1] if i + 1 < sss.size
      end
    end
    nil
  end

  def writing_duration
    return nil if lines.size < 2

    start_line = lines.min_by(&:created_at)
    showtime_lines = lines.reject { |line| line == start_line }
    first_line = showtime_lines.min_by(&:created_at)
    last_line = showtime_lines.max_by(&:created_at)
    return nil if last_line.nil? || first_line.nil?

    last_line.created_at - first_line.created_at
  end
end
