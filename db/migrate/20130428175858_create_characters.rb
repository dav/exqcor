class CreateCharacters < ActiveRecord::Migration
  def change
    create_table :characters do |t|
      t.string :name
      t.string :description
      t.string :play_id

      t.timestamps
    end
  end
end
