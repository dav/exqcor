require "test_helper"

class CharacterTest < ActiveSupport::TestCase
  test "only one vosd per script" do
    script = Script.create!(title: "Unique VOSD")
    script.characters.create!(name: "VOSD", role: "vosd")

    assert_db_constraint do
      script.characters.create!(name: "Another", role: "vosd")
    end

    assert_equal 1, script.characters.where(role: "vosd").count
  end

  test "cannot destroy vosd" do
    script = Script.create!(title: "Protect VOSD")
    vosd = script.ensure_vosd!

    vosd.destroy
    assert_not vosd.destroyed?
    assert Character.exists?(vosd.id)
  end
end
