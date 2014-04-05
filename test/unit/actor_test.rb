require 'test_helper'

class ActorTest < ActiveSupport::TestCase
  test "has characters" do
    assert actors(:oldman).characters.include? characters(:coxeter)
  end
end
