require 'test_helper'

class CharacterSectionTest < ActiveSupport::TestCase
  test "the fixtures" do
    assert_equal 3, CharacterSection.all.size

    play = plays(:play1)
    first_section = play.sections.first
    assert_equal 3, first_section.character_sections.size
  end
end
