require 'test_helper'

class CharacterTest < ActiveSupport::TestCase
  test "the truth" do
    coxeter = characters(:coxeter)
    assert_equal 'Donald', coxeter.name
    assert_equal actors(:oldman).id, coxeter.actor.id
  end
end
