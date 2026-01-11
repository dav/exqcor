require "test_helper"

class SectionTest < ActiveSupport::TestCase
  test "ensure_vosd does not duplicate character_sections" do
    script = Script.create!(title: "No Duplicates")
    section = script.sections.create!(name: "Scene 1")

    assert_equal 1, section.character_sections.count

    section.ensure_vosd!
    assert_equal 1, section.character_sections.count
  end
end
