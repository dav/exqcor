require 'test_helper'

class PlayTest < ActiveSupport::TestCase
  test "duplicate" do
    original_play_count = Play.count
    play = plays(:play1)
    new_play = play.duplicate(:title=> 'New Play')

    assert_equal(original_play_count+1, Play.count)

    assert_equal('New Play', new_play.title)
    assert_equal(play.description, new_play.description)

    assert_equal(play.characters.size, new_play.characters.size)
    assert_equal(actors(:oldman), new_play.characters[1].actor)

    assert_equal(play.sections.size, new_play.sections.size)

    play.sections.each_with_index do |section, i|
      new_section = new_play.sections[i]
      assert_equal(1, new_section.sub_sections.size)
      if i==0 # first section of fixture play
        new_first_ss = new_section.sub_sections.first
        new_first_line = new_first_ss.lines.first
        assert_equal(new_play.VOSD, new_first_line.character)
        assert_equal(lines(:p1s1ss1l1).text, new_first_line.text)

        assert_equal(3, section.character_sections.size)
        assert_equal(3, new_section.character_sections.size)
      end
    end
  end
end
