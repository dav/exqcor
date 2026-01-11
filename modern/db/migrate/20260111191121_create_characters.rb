class CreateCharacters < ActiveRecord::Migration[8.0]
  def change
    create_table :characters do |t|
      t.string :name
      t.string :description
      t.references :script, null: false, foreign_key: true
      t.references :actor, null: true, foreign_key: true

      t.timestamps
    end
  end
end
