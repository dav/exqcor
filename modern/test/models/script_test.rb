require "test_helper"

class ScriptTest < ActiveSupport::TestCase
  test "does not auto-create vosd on script create" do
    script = Script.create!(title: "No VOSD Yet")

    assert_nil script.VOSD
    assert_equal 0, script.characters.count
  end

  test "ensure_vosd is idempotent" do
    script = Script.create!(title: "Has VOSD")

    vosd = script.ensure_vosd!
    assert_equal "vosd", vosd.role
    assert_equal 1, script.characters.count

    assert_equal vosd.id, script.ensure_vosd!.id
    assert_equal 1, script.characters.count
  end

  test "vosd lookup uses role instead of name" do
    script = Script.create!(title: "Role Based")
    character = script.characters.create!(name: "Narrator", role: "vosd")

    assert_equal character.id, script.VOSD.id
  end
end
