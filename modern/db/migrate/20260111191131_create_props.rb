class CreateProps < ActiveRecord::Migration[8.0]
  def change
    create_table :props do |t|
      t.string :name
      t.string :description
      t.references :section, null: false, foreign_key: true

      t.timestamps
    end
  end
end
