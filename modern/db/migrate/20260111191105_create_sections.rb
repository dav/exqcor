class CreateSections < ActiveRecord::Migration[8.0]
  def change
    create_table :sections do |t|
      t.string :name
      t.string :description
      t.integer :ordering
      t.references :script, null: false, foreign_key: true

      t.timestamps
    end
  end
end
